package gameLogic_test

import (
	"fmt"
	"testing"

	"github.com/djpiper28/cards-against-humanity/gameLogic"
)

func TestNewCardDeckNoWhiteCards(t *testing.T) {
	whiteCards := make([]*gameLogic.WhiteCard, 0)
	blackCards := make([]*gameLogic.BlackCard, 2)

	_, err := gameLogic.NewCardDeck(whiteCards, blackCards)
	if err == nil {
		t.Log("Should not be able to make a deck with no white cards")
		t.FailNow()
	}
}

func TestNewCardDeckNoBlackCards(t *testing.T) {
	whiteCards := make([]*gameLogic.WhiteCard, 2)
	blackCards := make([]*gameLogic.BlackCard, 0)

	_, err := gameLogic.NewCardDeck(whiteCards, blackCards)
	if err == nil {
		t.Log("Should not be able to make a deck with no black cards")
		t.FailNow()
	}
}

func TestNewCardDeck(t *testing.T) {
	whiteCards := GetTestWhiteCards()
	blackCards := GetTestBlackCards()

	cardDeck, err := gameLogic.NewCardDeck(whiteCards, blackCards)
	if err != nil {
		t.Log(fmt.Sprintf("Error should be nil %s", err))
		t.FailNow()
	}

	if cardDeck == nil {
		t.Log("Card deck is nil")
		t.FailNow()
	}

	if len(cardDeck.WhiteCards) != len(whiteCards) {
		t.Log("White cards not set")
		t.FailNow()
	}

	if len(cardDeck.BlackCards) != len(blackCards) {
		t.Log("Black cards not set")
		t.FailNow()
	}
}

func TestGetNewWhiteCards(t *testing.T) {
	whiteCards := GetTestWhiteCards()
	testCardsLength := len(whiteCards)
	blackCards := GetTestBlackCards()

	cardDeck, err := gameLogic.NewCardDeck(whiteCards, blackCards)
	if err != nil {
		t.Log(fmt.Sprintf("Error should be nil %s", err))
		t.FailNow()
	}

	if cardDeck == nil {
		t.Log("Card deck is nil")
		t.FailNow()
	}

	var cardsToGet uint = 5
	cards, err := cardDeck.GetNewWhiteCards(uint(cardsToGet))
	if err != nil {
		t.Log(fmt.Sprintf("Error should be nil %s", err))
		t.FailNow()
	}

	if len(cards) != int(cardsToGet) {
		t.Log(fmt.Sprintf("Should have got %d cards, got %d", cardsToGet, len(cards)))
		t.FailNow()
	}

	expected := testCardsLength - int(cardsToGet)
	if len(cardDeck.WhiteCards) != expected {
		t.Log(fmt.Sprintf("Not enough cards where removed %d should be %d", len(cardDeck.WhiteCards), expected))
		t.FailNow()
	}
}