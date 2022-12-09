#[cfg(not(feature = "library"))]
use cosmwasm_std::entry_point;
use cosmwasm_std::{to_binary, Addr, Binary, Deps, DepsMut, Env, MessageInfo, Response, StdResult};
use cw2::set_contract_version;

use crate::error::ContractError;
use crate::msg::{ExecuteMsg, GetCardResponse, GetCardOwnerResponse, GetOwnersCardsResponse, IsOwnerResponse, InstantiateMsg, QueryMsg};
use crate::state::{Card, CardParams, CARDS, OWNER_CARDS, CARD_OWNER};
use crate::helpers::{ generate_random_by_level };

// version info for migration info
const CONTRACT_NAME: &str = "crates.io:card";
const CONTRACT_VERSION: &str = env!("CARGO_PKG_VERSION");

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn instantiate(
    deps: DepsMut,
    _env: Env,
    info: MessageInfo,
    msg: InstantiateMsg,
) -> Result<Response, ContractError> {
    // let state = State {
    //     count: msg.count,
    //     owner: info.sender.clone(),
    // };
    set_contract_version(deps.storage, CONTRACT_NAME, CONTRACT_VERSION)?;
    // STATE.save(deps.storage, &state)?;

    Ok(Response::new()
        .add_attribute("method", "instantiate")
        .add_attribute("owner", info.sender)
        .add_attribute("count", msg.count.to_string()))
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn execute(
    deps: DepsMut,
    _env: Env,
    info: MessageInfo,
    msg: ExecuteMsg,
) -> Result<Response, ContractError> {
    match msg {
        ExecuteMsg::Mint {} => execute::mint(deps, info),
        ExecuteMsg::Transfer {card_id, to} => execute::transfer(deps, info, card_id, to)
    }
}

pub mod execute {
    use cosmwasm_std::{Addr, StdError};

    use super::*;
    
    fn get_avatar_url() -> String{
        return "".to_string();
    }

    pub fn mint(deps: DepsMut, info: MessageInfo) -> Result<Response, ContractError> {
        let executor_addr: Addr = info.sender;
        let card_param = generate_random_by_level(0 as u8);
        let avatar_url: String = get_avatar_url();
        
        return save_card(deps, executor_addr, Card { param: card_param, avatar: avatar_url });
    }

    fn save_card(deps: DepsMut, executor: Addr, card: Card) -> Result<Response, ContractError>  {
        let card_id = card.get_key();
        let res = CARDS.save(deps.storage, card_id.clone(), &card);

        if res.is_err() {
            return Err(ContractError::CardExist{});
        }
        let res = CARD_OWNER.save(deps.storage, card_id.clone(), & executor);

        if res.is_err() {
            return Err(ContractError::UnknownErr{});
        }
        
        let put_card_to_owner_cards = |d: Option<Vec<String>> | -> StdResult<Vec<String>>   {
            let mut cards =d.unwrap_or_default();
            cards.push(card_id.clone());
            Ok(cards)
        };
        
        let res = OWNER_CARDS.update(deps.storage, executor.clone(), put_card_to_owner_cards);
        if res.is_err() {
            return Err(ContractError::UnknownErr{});
        }

        return Ok(Response::new().
        add_attribute("card_id", card_id).
        add_attribute("user", executor).
        add_attribute("action", "mint new card"));
    }


    pub fn transfer(deps: DepsMut, info: MessageInfo, card_id: String, to: Addr) -> Result<Response, ContractError> {
        
        let from = info.sender;
        let card = CARD_OWNER.load(deps.storage, card_id.clone());
        if card.is_err() {
            return Err(ContractError::CardNotExist{});
        }

        let owner = card.unwrap();
        if from.clone() != owner {
            return Err(ContractError::AccessDeny{});
        }

        _ = CARD_OWNER.update(deps.storage, card_id.clone(), |_: Option<Addr> | -> StdResult<Addr>   {Ok(to.clone())});

        let res = OWNER_CARDS.update(deps.storage, from.clone(), |d: Option<Vec<String>> | -> StdResult<Vec<String>>   {
            let mut cards = d.unwrap_or_default();
            let ff = cards.iter().enumerate().find(|&r| r.1.to_string() == card_id.clone());
            if ff.is_none() {
                return Err(StdError::NotFound { kind: "card not found".to_string() });
            } 
            cards.remove(ff.unwrap().0);
            Ok(cards)
        });
        if res.is_err() {
            return Err(ContractError::UnknownErr{});
        }

        _ = OWNER_CARDS.update(deps.storage, to.clone(), |d: Option<Vec<String>> | -> StdResult<Vec<String>>   {
            let mut cards = d.unwrap_or_default();
            cards.push(card_id.clone());
            Ok(cards)
        });


        Ok(Response::new()
        .add_attribute("from", from)
        .add_attribute("to", to)
        .add_attribute("action", "transfer"))
    }
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn query(deps: Deps, _env: Env, msg: QueryMsg) -> StdResult<Binary> {
    match msg {
        QueryMsg::GetCard { card_id } => to_binary(&query::get_card(deps, card_id)?),
        QueryMsg::GetCardOwner { card_id } => to_binary(&query::get_card_owner(deps, card_id)?),
        QueryMsg::GetOwnersCards { user_id } => to_binary(&query::get_owner_cards(deps, user_id)?),
        QueryMsg::IsOwner { card_id, user_id } => to_binary(&query::is_owner(deps, card_id, user_id)?),
    }
}

pub mod query {
    use super::*;
    
    pub fn get_card(deps: Deps, card_id: String) -> StdResult<GetCardResponse> {
        let card_info = CARDS.load(deps.storage, card_id).unwrap_or_default();

        Ok(GetCardResponse { card: card_info})
    }

    pub fn get_card_owner(deps: Deps, card_id: String) -> StdResult<GetCardOwnerResponse> {
        let card_owner = CARD_OWNER.load(deps.storage, card_id);
        if card_owner.is_err() {
            return Ok(GetCardOwnerResponse { user_id: Addr::unchecked("") });
        }
        Ok(GetCardOwnerResponse { user_id: card_owner.unwrap() })
    }
    
    pub fn get_owner_cards(deps: Deps, user_id: Addr) -> StdResult<GetOwnersCardsResponse> {
        let cards = OWNER_CARDS.load(deps.storage, user_id).unwrap_or_default();
        
        Ok(GetOwnersCardsResponse { card_ids: cards })
    }

    pub fn is_owner(deps: Deps, card_id: String, user_id: Addr) -> StdResult<IsOwnerResponse> {
        let card_owner = CARD_OWNER.load(deps.storage, card_id);
        if card_owner.is_err() {
            return Ok(IsOwnerResponse { is_owner: false });
        }
        Ok(IsOwnerResponse { is_owner: user_id == card_owner.unwrap() })
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use cosmwasm_std::testing::{mock_dependencies, mock_env, mock_info};
    use cosmwasm_std::{coins, from_binary};

    #[test]
    fn proper_initialization() {
        let mut deps = mock_dependencies();
        let empty_card: Card = Card::empty();

        let msg = InstantiateMsg { count: 17 };
        let info = mock_info("creator", &coins(1000, "earth"));

        // we can just call .unwrap() to assert this was a success
        let res = instantiate(deps.as_mut(), mock_env(), info, msg).unwrap();
        assert_eq!(0, res.messages.len());

        // it worked, let's query the state
        let res = query(deps.as_ref(), mock_env(), QueryMsg::GetCard { card_id: "".to_string() }).unwrap();
        let value: GetCardResponse = from_binary(&res).unwrap();
        assert_eq!(empty_card, value.card);
    }

    #[test]
    fn mint() {
        let mut deps = mock_dependencies();

        let msg = InstantiateMsg { count: 17 };
        let info = mock_info("creator", &coins(2, "token"));
        let _res = instantiate(deps.as_mut(), mock_env(), info, msg).unwrap();
        let sender: Addr = Addr::unchecked("sender");

        // beneficiary can release it
        let info = mock_info("sender", &coins(2, "token"));
        let msg = ExecuteMsg::Mint {};
        let _res = execute(deps.as_mut(), mock_env(), info, msg).unwrap();

        // should increase counter by 1
        let res = query(deps.as_ref(), mock_env(), QueryMsg::GetOwnersCards { user_id: sender.clone() }).unwrap();
        let value: GetOwnersCardsResponse = from_binary(&res).unwrap();
        assert_eq!(1, value.card_ids.len());

        let card_id_act = value.card_ids.first().unwrap();

        
        let res = query(deps.as_ref(), mock_env(), QueryMsg::GetCardOwner { card_id: card_id_act.to_string()}).unwrap();
        let value: GetCardResponse = from_binary(&res).unwrap();
        assert_ne!(CardParams::empty(), value.card.param);
    }

    #[test]
    fn test_move_card() {
        let mut deps = mock_dependencies();

        let user_from_str = "from_address";
        let user_from = Addr::unchecked(user_from_str);
        let user_to = Addr::unchecked("to_address");

        // Set owner of card
        
        let info = mock_info(user_from_str, &coins(2, "token"));
        let msg = ExecuteMsg::Mint {  } ;
        let res = execute(deps.as_mut(), mock_env(), info, msg);
        assert!(res.is_ok());
        let mut card_id_to_transfer = "".to_string();

        for attr in res.unwrap().attributes {
            if attr.key == "card_id".to_string() {
                card_id_to_transfer = attr.value;
                break;
            }
        }

        // Move card
        let info = mock_info(user_from_str, &coins(2, "token"));
        let msg = ExecuteMsg::Transfer { card_id: card_id_to_transfer.clone(), to: user_to.clone() } ;
        let _res = execute(deps.as_mut(), mock_env(), info, msg).unwrap();

        // Check that card has been transferred to the new owner
        let card = CARD_OWNER.load(deps.as_mut().storage, card_id_to_transfer.clone());
        assert_eq!(card.unwrap(), user_to.clone());

        // Check that the old owner's cards list no longer contains the card
        let owner_cards = OWNER_CARDS.load(deps.as_mut().storage, user_from.clone());
        assert!(!owner_cards.unwrap().contains(&card_id_to_transfer));

        // Check that the new owner's cards list contains the card
        let owner_cards = OWNER_CARDS.load(deps.as_mut().storage, user_to.clone());
        assert!(owner_cards.unwrap().contains(&card_id_to_transfer));
    }

}
