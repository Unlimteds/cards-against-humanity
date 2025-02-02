package gameLogic

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"

	"github.com/google/uuid"
)

type CardPack struct {
	Id         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	WhiteCards int       `json:"whiteCards"`
	BlackCards int       `json:"blackCards"`
	CardDeck   *CardDeck `json:"-"`
}

var AllPacks map[uuid.UUID]*CardPack
var AllWhiteCards []*WhiteCard
var AllBlackCards []*BlackCard

func GetWhiteCard(id int) (*WhiteCard, error) {
	if id < 0 || id >= len(AllWhiteCards) {
		return nil, errors.New("White card does not exist")
	}
	return AllWhiteCards[id], nil
}

func GetBlackCard(id int) (*BlackCard, error) {
	if id < 0 || id >= len(AllBlackCards) {
		return nil, errors.New("Black card does not exist")
	}
	return AllBlackCards[id], nil
}

func DefaultCardPack() *CardPack {
	for _, packValue := range AllPacks {
		return packValue
	}

	log.Println("Cannot find any packs for the default")
	return nil
}

func AccumalateCardPacks(packs []*CardPack) (*CardDeck, error) {
	if len(packs) == 0 {
		return nil, errors.New("At least one card pack must be selected")
	}

	decks := make([]*CardDeck, len(packs))
	for i, pack := range packs {
		decks[i] = pack.CardDeck
	}

	return AccumalateDecks(decks)
}

type cahJsonBlackCard struct {
	Text string `json:"text"`
	Pick int    `json:"pick"`
}

type cahJsonPack struct {
	Name             string `json:"name"`
	WhiteCardIndexes []int  `json:"white"`
	BlackCardIndexes []int  `json:"black"`
}

type cahJson struct {
	WhiteCards []string           `json:"white"`
	BlackCards []cahJsonBlackCard `json:"black"`
	Packs      []cahJsonPack      `json:"packs"`
}

const cahJsonFile = "packs/cah-all-compact.json"

func translateCahCards(data *cahJson) error {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		AllWhiteCards = make([]*WhiteCard, 0, len(data.WhiteCards))

		for i, cardText := range data.WhiteCards {
			AllWhiteCards = append(AllWhiteCards, NewWhiteCard(i, cardText))
		}
	}()

	AllBlackCards = make([]*BlackCard, 0, len(data.BlackCards))
	for i, blackCard := range data.BlackCards {
		AllBlackCards = append(AllBlackCards, NewBlackCard(i, blackCard.Text, uint(blackCard.Pick)))
	}

	wg.Wait()

	log.Printf("Found %d white cards and %d black cards", len(AllWhiteCards), len(AllBlackCards))
	return nil
}

func translateCahJson(data *cahJson) error {
	log.Println("Reading all cards")
	err := translateCahCards(data)
	if err != nil {
		log.Println("Cannot read the cards")
		return err
	}

	log.Println("Reading all packs")
	AllPacks = make(map[uuid.UUID]*CardPack)

	var wg sync.WaitGroup
	var lock sync.Mutex
	var terr error

	packs := 0
	for _, cahPack := range data.Packs {
		packs++
		wg.Add(1)
		go func(pack cahJsonPack) {
			defer wg.Done()

			id := uuid.New()
			whiteCards := make([]*WhiteCard, len(pack.WhiteCardIndexes))
			for i, whiteCardIndex := range pack.WhiteCardIndexes {
				whiteCards[i] = AllWhiteCards[whiteCardIndex]
			}

			blackCards := make([]*BlackCard, len(pack.BlackCardIndexes))
			for i, blackCardIndex := range pack.BlackCardIndexes {
				blackCards[i] = AllBlackCards[blackCardIndex]
			}

			deck, err := NewCardDeck(whiteCards, blackCards)
			if err != nil {
				log.Printf("Pack %s cannot be turned into a deck %s", pack.Name, err)
				lock.Lock()
				defer lock.Unlock()

				terr = err
				return
			}

			cardPack := CardPack{Id: id,
				CardDeck:   deck,
				Name:       pack.Name,
				WhiteCards: len(deck.WhiteCards),
				BlackCards: len(deck.BlackCards)}
			lock.Lock()
			defer lock.Unlock()
			AllPacks[id] = &cardPack
		}(cahPack)
	}

	wg.Wait()

	if terr != nil {
		log.Println("An error occurred whilst processing the decks (last error)", terr)
		AllPacks, AllWhiteCards, AllBlackCards = nil, nil, nil
	}

	log.Printf("Created %d packs of cards", packs)
	return terr
}

func LoadPacks() error {
	if AllPacks != nil {
		log.Println("Data is already loaded")
		return nil
	}

	log.Println("Reading data file", cahJsonFile)

	dataFileContents, err := os.ReadFile(cahJsonFile)
	if err != nil {
		log.Println("Cannot read data file", cahJsonFile, err)
		return err
	}

	log.Println("Parsing data file")

	var cahData cahJson
	err = json.Unmarshal(dataFileContents, &cahData)
	if err != nil {
		log.Println("Cannot parse data file", err)
		return err
	}

	err = translateCahJson(&cahData)
	if err != nil {
		log.Println("Cannot translate the data file to the internal struct")
		return err
	}
	return nil
}
