package contracts

import (
	"gameoflife/utils"
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
	cards := NewCards(1000)

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
