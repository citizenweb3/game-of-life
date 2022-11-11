package contracts

import (
	"gameoflife/system/mocks"
	"gameoflife/utils"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func generateDummyCard(id string, num int) Card {
	return Card{
		Id: utils.CardID(id),
		Params: CardParams{
			Hp:       uint64(100 + num),
			Level:    uint8(200 + num),
			Strength: uint64(300 + num),
			Accuracy: uint64(400 + num),
		},
	}
}

func TestCards_GetOwnersCards(t *testing.T) {
	owner1List := []utils.CardID{"11111", "22222", "44444"}
	cardList := map[utils.CardID]Card{
		"11111": generateDummyCard("11111", 1),
		"22222": generateDummyCard("22222", 2),
		"33333": generateDummyCard("33333", 3),
		"44444": generateDummyCard("44444", 4),
	}
	cardsContract := Cards{
		ownerCards: map[utils.UserID][]utils.CardID{
			"owner1": owner1List,
		},
		cards: cardList,
	}
	expectedO1 := []Card{
		cardList["11111"],
		cardList["22222"],
		cardList["44444"],
	}
	actualCardListRealUser := cardsContract.GetOwnersCards("owner1")
	actualCardListNotExistUser := cardsContract.GetOwnersCards("owner_unexpected")

	assert.Equal(t, expectedO1, actualCardListRealUser)
	assert.Empty(t, actualCardListNotExistUser)
}

func TestCards_GetCardOwner(t *testing.T) {
	cardsContract := Cards{
		cardOwner: map[utils.CardID]utils.UserID{
			"11111": "owner1",
			"22222": "owner2",
		},
	}
	actualCardOwner, err1 := cardsContract.GetCardOwner("11111")
	actualCardOwnerNotExistCard, err2 := cardsContract.GetCardOwner("sddskl")

	require.Nil(t, err1)
	assert.Equal(t, utils.UserID("owner1"), actualCardOwner)

	require.Error(t, err2)
	assert.Equal(t, utils.ErrCardNotExist, err2)
	assert.Equal(t, utils.UserID(""), actualCardOwnerNotExistCard)
}

func TestCards_GetCardProperties(t *testing.T) {
	cardList := map[utils.CardID]Card{
		"11111": generateDummyCard("11111", 1),
		"22222": generateDummyCard("22222", 2),
	}
	cardsContract := Cards{
		cards: cardList,
	}
	expectedCardParam := cardList["11111"].Params
	actualCardParam, err1 := cardsContract.GetCardProperties("11111")
	actualCardParamNotExistCard, err2 := cardsContract.GetCardProperties("77777")

	require.NoError(t, err1)
	assert.Equal(t, expectedCardParam, actualCardParam)

	require.Error(t, err2)
	assert.Equal(t, utils.ErrCardNotExist, err2)
	assert.Equal(t, CardParams{}, actualCardParamNotExistCard)
}

func TestCards_MintNewCard(t *testing.T) {
	cards := NewCards(nil, 1000)

	cid11, err11 := cards.MintNewCard("user1")
	cid12, err12 := cards.MintNewCard("user1")
	cid13, err13 := cards.MintNewCard("user1")
	cid21, err21 := cards.MintNewCard("user2")

	require.NoError(t, err11)
	require.NoError(t, err12)
	require.NoError(t, err13)
	require.NoError(t, err21)

	expOwnerCards := map[utils.UserID][]utils.CardID{
		"user1": {cid11, cid12, cid13},
		"user2": {cid21},
	}
	require.Equal(t, len(expOwnerCards), len(cards.ownerCards))
	for userID, cardIDs := range expOwnerCards {
		assert.Equal(t, cardIDs, cards.ownerCards[userID])
	}

	expCardOwner := map[utils.CardID]utils.UserID{
		cid11: "user1",
		cid12: "user1",
		cid13: "user1",
		cid21: "user2",
	}
	require.Equal(t, len(expCardOwner), len(cards.cardOwner))
	for cardID, userID := range expCardOwner {
		assert.Equal(t, userID, cards.cardOwner[cardID])
	}

	assert.Equal(t, 4, len(cards.cards))

	assert.NotEqualValues(t, cid11, cid12)
	assert.NotEqualValues(t, cid11, cid13)
	assert.NotEqualValues(t, cid11, cid21)
	assert.NotEqualValues(t, cid12, cid13)
	assert.NotEqualValues(t, cid12, cid21)
	assert.NotEqualValues(t, cid13, cid21)
}

func TestCards_IsOwner(t *testing.T) {
	cards := Cards{
		cardOwner: map[utils.CardID]utils.UserID{"11111": "user1"},
	}
	actCorrect := cards.IsOwner("11111", "user1")
	actCardNotExits := cards.IsOwner("22222", "user1")
	actIsNotOwner := cards.IsOwner("11111", "user2")

	assert.NoError(t, actCorrect)
	require.Error(t, actCardNotExits)
	assert.Equal(t, utils.ErrCardNotExist, actCardNotExits)
	require.Error(t, actIsNotOwner)
	assert.Equal(t, utils.ErrPermDenie, actIsNotOwner)
}

func TestCards_Transfer(t *testing.T) {
	t.Run("justOneCard", func(t *testing.T) {
		cards := Cards{
			cardOwner:  map[utils.CardID]utils.UserID{"11111": "user1"},
			ownerCards: map[utils.UserID][]utils.CardID{"user1": {"11111"}},
		}
		expCardOwner := map[utils.CardID]utils.UserID{"11111": "user2"}
		expOwnerCards := map[utils.UserID][]utils.CardID{"user2": {"11111"}, "user1": {}}

		err := cards.Transfer("11111", "user1", "user2")

		require.NoError(t, err)
		assert.Equal(t, expCardOwner, cards.cardOwner)
		assert.Equal(t, expOwnerCards, cards.ownerCards)
	})

	t.Run("moreThanOneCard", func(t *testing.T) {
		cards := Cards{
			cardOwner:  map[utils.CardID]utils.UserID{"11111": "user1", "22222": "user1", "33333": "user2"},
			ownerCards: map[utils.UserID][]utils.CardID{"user1": {"11111", "22222"}, "user2": {"33333"}},
		}
		expCardOwner := map[utils.CardID]utils.UserID{"11111": "user2", "22222": "user1", "33333": "user2"}
		expOwnerCards := map[utils.UserID][]utils.CardID{"user2": {"11111", "33333"}, "user1": {"22222"}}

		err := cards.Transfer("11111", "user1", "user2")

		require.NoError(t, err)
		assert.Equal(t, expCardOwner, cards.cardOwner)
		assert.Equal(t, len(expOwnerCards["user2"]), len(cards.ownerCards["user2"]))
		for _, expCard := range expOwnerCards["user2"] {
			assert.Contains(t, cards.ownerCards["user2"], expCard)
		}
		assert.Equal(t, len(expOwnerCards["user1"]), len(cards.ownerCards["user1"]))
		for _, expCard := range expOwnerCards["user1"] {
			assert.Contains(t, cards.ownerCards["user1"], expCard)
		}
	})

	t.Run("is not owner", func(t *testing.T) {
		cards := Cards{
			cardOwner:  map[utils.CardID]utils.UserID{"11111": "user1"},
			ownerCards: map[utils.UserID][]utils.CardID{"user1": {"11111"}},
		}
		err := cards.Transfer("11111", "user2", "user3")
		assert.Error(t, err)
		assert.Equal(t, utils.ErrPermDenie, err)
	})

	t.Run("is not exist", func(t *testing.T) {
		cards := Cards{
			cardOwner:  map[utils.CardID]utils.UserID{"11111": "user1"},
			ownerCards: map[utils.UserID][]utils.CardID{"user1": {"11111"}},
		}
		err := cards.Transfer("00000", "user1", "user2")
		assert.Error(t, err)
		assert.Equal(t, utils.ErrCardNotExist, err)
	})
}

func Test_getRandBeetweenUint64(t *testing.T) {
	caseCount := 10

	for i := 0; i < caseCount; i++ {
		p1 := rand.Uint64() % 1000000
		p2 := rand.Uint64() % 1000000
		maxBase := p1
		minBase := p2
		if p1 < p2 {
			maxBase, minBase = minBase, maxBase
		}
		maxBase += MIXED_ADD

		res := getRandBeetweenUint64(p1, p2)
		assert.LessOrEqual(t, minBase, res)
		assert.LessOrEqual(t, res, maxBase)
	}
}

func TestCards_mixedCards(t *testing.T) {
	card1 := generateDummyCard("1", 1)
	card2 := generateDummyCard("2", 2)
	user1 := utils.UserID("user1")
	cardsContract := Cards{
		cardOwner: map[utils.CardID]utils.UserID{
			card1.Id: user1,
			card2.Id: user1,
		},
		ownerCards: map[utils.UserID][]utils.CardID{
			user1: {card1.Id, card2.Id},
		},
		cards: map[utils.CardID]Card{
			card1.Id: card1,
			card2.Id: card2,
		},
	}

	card := cardsContract.mixedCards(card1.Id, card2.Id)
	assert.LessOrEqual(t, card1.Params.Hp, card.Params.Hp)
	assert.LessOrEqual(t, card.Params.Hp, card2.Params.Hp+MIXED_ADD)

	assert.LessOrEqual(t, card1.Params.Accuracy, card.Params.Accuracy)
	assert.LessOrEqual(t, card.Params.Accuracy, card2.Params.Accuracy+MIXED_ADD)

	assert.LessOrEqual(t, card1.Params.Strength, card.Params.Strength)
	assert.LessOrEqual(t, card.Params.Strength, card2.Params.Strength+MIXED_ADD)

	assert.Equal(t, card2.Params.Level, card.Params.Level)
}

func TestCards_UnFreeze(t *testing.T) {
	card1 := utils.CardID("card1")
	card2 := utils.CardID("card2")
	card1_1 := utils.CardID("card1_1")
	card2_1 := utils.CardID("card2_1")
	card3 := utils.CardID("card3")
	card4 := utils.CardID("card4")
	card5 := utils.CardID("card5")
	cardUnexist := utils.CardID("unexistCardId")
	user1 := utils.UserID("user1")
	user2 := utils.UserID("user2")
	user3 := utils.UserID("user3")
	systemMocked := mocks.NewSystemI(t)
	systemMocked.On("GetCurrentTime").Return(int64(10000))
	cardsContract := Cards{
		system: systemMocked,
		cardOwner: map[utils.CardID]utils.UserID{
			card1:   user1,
			card2:   user1,
			card1_1: user3,
			card2_1: user3,
			card3:   user2,
			card4:   user2,
			card5:   user2,
		},
		ownerCards: map[utils.UserID][]utils.CardID{
			user1: {card1, card2},
			user2: {card3, card4, card5},
			user3: {card1_1, card2_1},
		},
		cards: map[utils.CardID]Card{
			card1:   generateDummyCard("1", 1),
			card2:   generateDummyCard("2", 2),
			card3:   generateDummyCard("3", 3),
			card4:   generateDummyCard("4", 4),
			card5:   generateDummyCard("5", 5),
			card1_1: generateDummyCard("1_1", 11),
			card2_1: generateDummyCard("2_1", 21),
		},
		freezedPair: map[utils.CardID]utils.CardID{
			card1:   card2,
			card2:   card1,
			card1_1: card2_1,
			card2_1: card1_1,
			card3:   card4,
			card4:   card3,
		},
		freezedUntil: map[string]int64{
			card1.ToString():   9000,
			card1_1.ToString(): 9000,
			card3.ToString():   11000,
		},
	}
	t.Run("Correct card1", func(t *testing.T) {
		err := cardsContract.UnFreeze(card1, user1)
		require.NoError(t, err)
		// check delete card1, card2
		require.Equal(t, len(cardsContract.ownerCards[user1]), 1)

		assert.NotContains(t, cardsContract.ownerCards[user1], card1)
		assert.NotContains(t, cardsContract.ownerCards[user1], card2)

		assert.NotContains(t, cardsContract.freezedPair, card1)
		assert.NotContains(t, cardsContract.freezedPair, card2)

		assert.NotContains(t, cardsContract.freezedUntil, card1)
		assert.NotContains(t, cardsContract.freezedUntil, card2)

		assert.NotContains(t, cardsContract.cards, card1)
		assert.NotContains(t, cardsContract.cards, card2)

		assert.NotContains(t, cardsContract.cardOwner, card1)
		assert.NotContains(t, cardsContract.cardOwner, card2)

		// check existing new card
		cardNew := cardsContract.ownerCards[user1][0]

		assert.Contains(t, cardsContract.cardOwner, cardNew)
		assert.Contains(t, cardsContract.cards, cardNew)
	})

	t.Run("Correct card2", func(t *testing.T) {
		err := cardsContract.UnFreeze(card2_1, user3)
		require.NoError(t, err)
		// check delete card1, card2
		require.Equal(t, len(cardsContract.ownerCards[user3]), 1)

		assert.NotContains(t, cardsContract.ownerCards[user3], card1_1)
		assert.NotContains(t, cardsContract.ownerCards[user3], card2_1)

		assert.NotContains(t, cardsContract.freezedPair, card1_1)
		assert.NotContains(t, cardsContract.freezedPair, card2_1)

		assert.NotContains(t, cardsContract.freezedUntil, card1_1)
		assert.NotContains(t, cardsContract.freezedUntil, card2_1)

		assert.NotContains(t, cardsContract.cards, card1_1)
		assert.NotContains(t, cardsContract.cards, card2_1)

		assert.NotContains(t, cardsContract.cardOwner, card1_1)
		assert.NotContains(t, cardsContract.cardOwner, card2_1)

		// check existing new card
		cardNew := cardsContract.ownerCards[user3][0]

		assert.Contains(t, cardsContract.cardOwner, cardNew)
		assert.Contains(t, cardsContract.cards, cardNew)
	})

	t.Run("Card not exist", func(t *testing.T) {
		err := cardsContract.UnFreeze(cardUnexist, user1)

		assert.ErrorIs(t, err, utils.ErrCardNotExist)
	})
	t.Run("Is not owner", func(t *testing.T) {
		err := cardsContract.UnFreeze(card3, user1)

		assert.ErrorIs(t, err, utils.ErrPermDenie)
	})

	t.Run("Cards is still freezed", func(t *testing.T) {
		err := cardsContract.UnFreeze(card3, user2)
		assert.ErrorIs(t, err, utils.ErrCardFreezed)
		err = cardsContract.UnFreeze(card4, user2)
		assert.ErrorIs(t, err, utils.ErrCardFreezed)
	})
	t.Run("Cards is not freezed", func(t *testing.T) {
		err := cardsContract.UnFreeze(card5, user2)
		assert.ErrorIs(t, err, utils.ErrCardNotFreezed)
	})
}

func TestCards_Freeze(t *testing.T) {

	card1 := utils.CardID("card1")
	card2 := utils.CardID("card2")
	card3 := utils.CardID("card3")
	card4 := utils.CardID("card4")
	card5 := utils.CardID("card5")
	card6 := utils.CardID("card6")
	cardUnexist := utils.CardID("unexistCardId")
	user1 := utils.UserID("user1")
	user2 := utils.UserID("user2")
	currentSysTime := int64(1000)
	freezeTime := int64(1000)
	systemMocked := mocks.NewSystemI(t)
	systemMocked.On("GetCurrentTime").Return(currentSysTime)
	cardsContract := Cards{
		system:     systemMocked,
		freezeTime: freezeTime,
		cardOwner: map[utils.CardID]utils.UserID{
			card1: user1,
			card2: user1,
			card3: user1,
			card4: user1,
			card5: user1,
			card6: user2,
		},
		ownerCards: map[utils.UserID][]utils.CardID{
			user1: {card1, card2, card3, card4, card5},
			user2: {card6},
		},
		cards: map[utils.CardID]Card{
			card1: generateDummyCard("1", 1),
			card2: generateDummyCard("2", 2),
			card3: generateDummyCard("3", 3),
			card4: generateDummyCard("4", 4),
			card5: generateDummyCard("5", 5),
			card5: generateDummyCard("6", 6),
		},
		freezedPair: map[utils.CardID]utils.CardID{
			card1: card2,
			card2: card1,
		},
		freezedUntil: map[string]int64{
			card1.ToString(): currentSysTime - 20,
		},
	}

	t.Run("Correct", func(t *testing.T) {
		err := cardsContract.Freeze(card3, card4, user1)
		require.NoError(t, err)

		assert.True(t, cardsContract.IsFreezed(card3))
		assert.True(t, cardsContract.IsFreezed(card4))

		require.Contains(t, cardsContract.freezedPair, card3)
		assert.Contains(t, cardsContract.freezedPair[card3], card4)
		require.Contains(t, cardsContract.freezedPair, card4)
		assert.Contains(t, cardsContract.freezedPair[card4], card3)

		require.Contains(t, cardsContract.freezedUntil, card3.ToString())
		assert.NotContains(t, cardsContract.freezedUntil, card4.ToString())
		assert.Equal(t, freezeTime+currentSysTime, cardsContract.freezedUntil[card3.ToString()])

	})

	t.Run("Card is equal", func(t *testing.T) {
		err := cardsContract.Freeze(card5, card5, user1)
		assert.ErrorIs(t, err, utils.ErrCardsEqual)
	})

	t.Run("Card is not exist", func(t *testing.T) {
		err := cardsContract.Freeze(card5, cardUnexist, user1)
		assert.ErrorIs(t, err, utils.ErrCardNotExist)
	})

	t.Run("It's not owner one of cards", func(t *testing.T) {
		err := cardsContract.Freeze(card5, card6, user1)
		assert.ErrorIs(t, err, utils.ErrPermDenie)

		err = cardsContract.Freeze(card6, card5, user1)
		assert.ErrorIs(t, err, utils.ErrPermDenie)
	})

	t.Run("Card is already freezed", func(t *testing.T) {
		err := cardsContract.Freeze(card1, card5, user1)
		assert.ErrorIs(t, err, utils.ErrCardFreezed)
		err = cardsContract.Freeze(card5, card1, user1)
		assert.ErrorIs(t, err, utils.ErrCardFreezed)
	})
}
