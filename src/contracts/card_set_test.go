package contracts_test

import (
	"errors"
	"fmt"
	"gameoflife/contracts"
	"gameoflife/contracts/mocks"
	"gameoflife/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCardSet_AddCardToSet(t *testing.T) {

	cardMock := mocks.NewCardsI(t)
	cardMock.On("IsOwner", utils.CardID("card1"), utils.UserID("user1")).Return(nil)
	cardMock.On("IsOwner", utils.CardID("card2"), utils.UserID("user1")).Return(nil)
	cardMock.On("IsOwner", utils.CardID("card1"), utils.UserID("user2")).Return(utils.ErrPermDenie)

	t.Run("Correct", func(t *testing.T) {
		cardSetCorrect := contracts.NewCardSet(cardMock, 2)
		expectedSet := []utils.CardID{"card1"}

		err := cardSetCorrect.AddCardToSet("user1", "card1")

		require.NoError(t, err)

		actualSet := cardSetCorrect.GetActualSet("user1")
		assert.Equal(t, expectedSet, actualSet)

	})

	cardSet := contracts.NewCardSet(cardMock, 2)
	cardSet.AddCardToSet("user1", "card1")

	t.Run("Is not owner", func(t *testing.T) {
		err := cardSet.AddCardToSet("user2", "card1")
		require.Error(t, err)
		assert.True(t, errors.Is(err, utils.ErrPermDenie))
	})

	t.Run("Set is alredy full", func(t *testing.T) {
		cardSet := contracts.NewCardSet(cardMock, 1)
		err := cardSet.AddCardToSet("user1", "card1")
		require.NoError(t, err)
		err = cardSet.AddCardToSet("user1", "card2")
		require.Error(t, err)
		assert.True(t, errors.Is(err, utils.ErrSetTooMuchCards))
	})

	t.Run("Card is alredy in set", func(t *testing.T) {
		err := cardSet.AddCardToSet("user1", "card1")
		require.Error(t, err)
		assert.True(t, errors.Is(err, utils.ErrSetAlreadyInSet))
	})
}

func TestCardSet_RemoveCardFromSet(t *testing.T) {
	cardMock := mocks.NewCardsI(t)
	cardMock.On("IsOwner", utils.CardID("card1"), utils.UserID("user1")).Return(nil).Once()
	cardMock.On("IsOwner", utils.CardID("card2"), utils.UserID("user1")).Return(nil).Once()

	cardSet := contracts.NewCardSet(cardMock, 2)
	err := cardSet.AddCardToSet("user1", "card1")
	require.NoError(t, err, "Error in add card to set")
	err = cardSet.AddCardToSet("user1", "card2")
	require.NoError(t, err, "Error in add card to set")

	t.Run("Correct", func(t *testing.T) {
		expectedSet := []utils.CardID{"card2"}
		err := cardSet.RemoveCardFromSet("user1", "card1")
		assert.NoError(t, err)
		actualSet := cardSet.GetActualSet("user1")
		assert.Equal(t, expectedSet, actualSet)

	})

	t.Run("Empty set", func(t *testing.T) {
		err := cardSet.RemoveCardFromSet("user2", "card1")
		assert.Error(t, err)
		assert.True(t, errors.Is(err, utils.ErrSetIsEmpty))
	})

	t.Run("Card is not in set", func(t *testing.T) {
		err := cardSet.RemoveCardFromSet("user1", "card3")
		assert.Error(t, err)
		assert.True(t, errors.Is(err, utils.ErrSetCardIsNotInSet))
	})
}

func TestCardSet_ChangeCardFromSet(t *testing.T) {

	user1 := utils.UserID("user1")
	user2 := utils.UserID("user2")

	card0 := utils.CardID("card0")
	card1 := utils.CardID("card1")
	card2 := utils.CardID("card2")
	card3 := utils.CardID("card3")
	card4 := utils.CardID("card4")
	card5 := utils.CardID("card5")
	cardUser2 := utils.CardID("cardUser2")
	unexistingCard := utils.CardID("unexistingCard")

	cardMock := mocks.NewCardsI(t)

	cardMock.On("IsOwner", card0, user1).Return(nil).Once()
	cardMock.On("IsOwner", card1, user1).Return(nil).Once()
	cardMock.On("IsOwner", card2, user1).Return(nil).Twice()
	cardMock.On("IsOwner", card3, user1).Return(nil).Once()
	cardMock.On("IsOwner", card5, user1).Return(nil).Once()

	cardMock.On("IsOwner", cardUser2, user1).Return(utils.ErrPermDenie).Once()
	cardMock.On("IsOwner", unexistingCard, user1).Return(utils.ErrCardNotExist).Once()

	cardMock.On("IsOwner", cardUser2, user2).Return(nil).Once()

	cardSet := contracts.NewCardSet(cardMock, 3)
	err := cardSet.AddCardToSet(user1, card0)
	require.NoError(t, err, "Error in add card to set")
	err = cardSet.AddCardToSet(user1, card1)
	require.NoError(t, err, "Error in add card to set")
	err = cardSet.AddCardToSet(user1, card2)
	require.NoError(t, err, "Error in add card to set")

	t.Run("Correct", func(t *testing.T) {
		err := cardSet.ChangeCardFromSet(user1, card0, card3)

		assert.NoError(t, err)

		actualSet := cardSet.GetActualSet(user1)

		assert.Contains(t, actualSet, card3)
		assert.NotContains(t, actualSet, card0)
	})

	t.Run("Card not exist", func(t *testing.T) {
		err := cardSet.ChangeCardFromSet(user1, card1, unexistingCard)
		assert.ErrorIs(t, err, utils.ErrCardNotExist)
	})

	t.Run("Wrong card owner", func(t *testing.T) {
		err := cardSet.ChangeCardFromSet(user1, card1, cardUser2)
		assert.ErrorIs(t, err, utils.ErrPermDenie)
	})

	t.Run("Card already in set", func(t *testing.T) {
		err := cardSet.ChangeCardFromSet(user1, card1, card2)
		assert.ErrorIs(t, err, utils.ErrSetAlreadyInSet)
	})

	t.Run("Card in not from set", func(t *testing.T) {
		err := cardSet.ChangeCardFromSet(user1, card4, card5)
		assert.ErrorIs(t, err, utils.ErrSetCardIsNotInSet)
	})

	t.Run("Has no set", func(t *testing.T) {
		err := cardSet.ChangeCardFromSet(user2, card1, cardUser2)
		assert.ErrorIs(t, err, utils.ErrSetIsEmpty)
	})
}

func TestCardSet_SetUserAttributes(t *testing.T) {
	maxCardNum := 10
	cs := contracts.NewCardSet(nil, uint8(maxCardNum))
	pushedParam := contracts.CardParams{Hp: 10, Level: 11, Strength: 12, Accuracy: 13}
	pushedParam2 := contracts.CardParams{Hp: 20, Level: 21, Strength: 22, Accuracy: 23}
	emptyParam := contracts.CardParams{Hp: 0, Level: 0, Strength: 0, Accuracy: 0}

	// push in empty set
	err := cs.SetUserAttribute("user1", 1, pushedParam)

	assert.NoError(t, err)
	actAttr := cs.GetUserAttributes("user1")
	assert.Equal(t, maxCardNum, len(actAttr))
	assert.Equal(t, pushedParam, actAttr[1])
	assert.Equal(t, emptyParam, actAttr[0], "attribute 0 is not empty")
	for i := 2; i < maxCardNum; i++ {
		assert.Equal(t, emptyParam, actAttr[i], fmt.Sprintf("attribute %d is not empty", i))
	}

	t.Run("Correct in Exist set", func(t *testing.T) {
		err := cs.SetUserAttribute("user1", 2, pushedParam2)
		assert.NoError(t, err)
		actAttr := cs.GetUserAttributes("user1")

		assert.Equal(t, maxCardNum, len(actAttr))
		// don't check actAttr[1] value because it may be update in case "update value" working in parallel
		assert.NotEqual(t, emptyParam, actAttr[1])
		assert.Equal(t, pushedParam2, actAttr[2])
		assert.Equal(t, emptyParam, actAttr[0], "attribute 0 is not empty")
		for i := 3; i < maxCardNum; i++ {
			assert.Equal(t, emptyParam, actAttr[i], fmt.Sprintf("attribute %d is not empty", i))
		}
	})

	t.Run("update value", func(t *testing.T) {
		err := cs.SetUserAttribute("user1", 1, pushedParam2)
		assert.NoError(t, err)
		actAttr := cs.GetUserAttributes("user1")

		assert.Equal(t, maxCardNum, len(actAttr))
		assert.Equal(t, pushedParam2, actAttr[1])
		assert.Equal(t, emptyParam, actAttr[0], "attribute 0 is not empty")
		// ignore element num 2 because "Correct in Exist set" working in parallel
		for i := 3; i < maxCardNum; i++ {
			assert.Equal(t, emptyParam, actAttr[i], fmt.Sprintf("attribute %d is not empty", i))
		}
	})

	t.Run("OutOfRange", func(t *testing.T) {
		err := cs.SetUserAttribute("user1", uint8(maxCardNum+1), pushedParam2)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, utils.ErrOutOfSetRange))
	})
}