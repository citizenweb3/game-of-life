#[cfg(not(feature = "library"))]
use cosmwasm_std::entry_point;
use cosmwasm_std::{to_binary, Binary, Deps, DepsMut, Env, MessageInfo, Response, StdResult};
use cosmwasm_std::Addr;
use cw2::set_contract_version;



use crate::error::ContractError;
use crate::msg::{ExecuteMsg, InstantiateMsg, QueryMsg, GetActualSetResponse};
use crate::state::{CardsSet, CARDS_SETS, COUNT_CARDS_IN_SET};

// version info for migration info
const CONTRACT_NAME: &str = "crates.io:cards-set";
const CONTRACT_VERSION: &str = env!("CARGO_PKG_VERSION");

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn instantiate(
    deps: DepsMut,
    _env: Env,
    info: MessageInfo,
    _: InstantiateMsg,
) -> Result<Response, ContractError> {
    // let state = State {
    //     count: msg.count,
    //     owner: info.sender.clone(),
    // };
    set_contract_version(deps.storage, CONTRACT_NAME, CONTRACT_VERSION)?;
    // STATE.save(deps.storage, &state)?;

    Ok(Response::new()
        .add_attribute("method", "instantiate")
        // .add_attribute("count", msg.count.to_string())
        .add_attribute("owner", info.sender))
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn execute(
    deps: DepsMut,
    _env: Env,
    info: MessageInfo,
    msg: ExecuteMsg,
) -> Result<Response, ContractError> {
    match msg {
        ExecuteMsg::SetCardToSet {num_in_set, card_id} => execute::set_card_to_set(deps, info, num_in_set, card_id),
        ExecuteMsg::SetUserAttribute { num_in_set, value } => execute::set_user_attribute(deps, info, num_in_set, value),
    }
}

pub mod execute {
    //use schemars::_serde_json::Error;

    use crate::state::UserParams;

    use super::*;

    pub fn set_card_to_set(deps: DepsMut, info: MessageInfo, num_in_set: usize, card_id: String) -> Result<Response, ContractError> {
        if num_in_set > COUNT_CARDS_IN_SET {
            return Err(ContractError::TooMuchCardNum {});
        }

        // toDo check IsOwner(cardId, info.sender)
        // return Err(ContractError::Unauthorized {});

        let put_card_set = |d: Option<CardsSet> | -> StdResult<CardsSet>   {
            let mut cards_set =d.unwrap_or_default();
            cards_set.cards_id[num_in_set] = card_id;
            return Ok(cards_set)
        };

        CARDS_SETS.update(deps.storage, info.sender, put_card_set)?;
        Ok(Response::new().
        add_attribute("action", "setCardToSet"))
    }

    pub fn set_user_attribute(deps: DepsMut, info: MessageInfo, num_in_set: usize, value: UserParams) -> Result<Response, ContractError> {
        if num_in_set > COUNT_CARDS_IN_SET {
            return Err(ContractError::TooMuchCardNum {});
        }

        let put_card_set = |d: Option<CardsSet> | -> StdResult<CardsSet>   {
            let mut cards_set = d.unwrap_or_default();
            cards_set.user_influence[num_in_set] = value;
            // to do add check to more than 100 total value for any param
            // if ... {
            //      return Err(ContractError::TooMuchInfluense)
            // }
            return Ok(cards_set)
        };

        CARDS_SETS.update(deps.storage, info.sender, put_card_set)?;
        Ok(Response::new().add_attribute("setAttribute", "Done"))
    }
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn query(deps: Deps, _env: Env, msg: QueryMsg) -> StdResult<Binary> {
    match msg {
        QueryMsg::GetActualSet {user} => to_binary(&query::get_actual_set(deps, user)?),
    }
}

pub mod query {
    use super::*;

    pub fn get_actual_set(deps: Deps, user: Addr) -> StdResult<GetActualSetResponse> {
        let set_val = CARDS_SETS.load(deps.storage, user).unwrap_or_default();
        Ok(GetActualSetResponse { set: set_val })
    }
}

#[cfg(test)]
mod tests {
    use crate::state::UserParams;

    use super::*;
    use cosmwasm_std::testing::{mock_dependencies, mock_env, mock_info};
    use cosmwasm_std::{coins, from_binary};

    #[test]
    fn proper_initialization() {
        let mut deps = mock_dependencies();
        let def_user_param: UserParams= UserParams { ampers: 0, volts: 0 };

        let empty_set: CardsSet = CardsSet {
             cards_id: ["".to_string(), "".to_string(), "".to_string(), "".to_string(), "".to_string()], 
             user_influence: [def_user_param, def_user_param, def_user_param, def_user_param, def_user_param],
            };

        let msg = InstantiateMsg { };
        let info = mock_info("creator", &coins(1000, "earth"));

        // we can just call .unwrap() to assert this was a success
        let res = instantiate(deps.as_mut(), mock_env(), info, msg).unwrap();
        assert_eq!(0, res.messages.len());
        // it worked, let's query the state
        let res = query(deps.as_ref(), mock_env(), QueryMsg::GetActualSet { user: Addr::unchecked("anyone")}).unwrap();
        let value: GetActualSetResponse = from_binary(&res).unwrap();
        assert_eq!(empty_set, value.set)

    }

    #[test]
    fn set_card_to_set() {
        let mut deps = mock_dependencies();

        let msg = InstantiateMsg {};
        let info = mock_info("creator", &coins(2, "token"));
        let _res = instantiate(deps.as_mut(), mock_env(), info, msg).unwrap();

        // beneficiary can release it
        let info = mock_info("anyone", &coins(2, "token"));
        let msg = ExecuteMsg::SetCardToSet { num_in_set: 0, card_id: "11111".to_string() };
        let _res = execute(deps.as_mut(), mock_env(), info, msg).unwrap();

        // create new card set
        let addr =Addr::unchecked("anyone");
        let res = query(deps.as_ref(), mock_env(), QueryMsg::GetActualSet { user: addr }).unwrap();
        let value: GetActualSetResponse = from_binary(&res).unwrap();
        assert_eq!("11111".to_string(), value.set.cards_id[0]);
        assert_eq!("".to_string(), value.set.cards_id[1]);

        // // update new card set 
        let addr =Addr::unchecked("anyone");
        let info = mock_info("anyone", &coins(2, "token"));
        let msg = ExecuteMsg::SetCardToSet { num_in_set: 0, card_id: "22222".to_string() };
        let _res = execute(deps.as_mut(), mock_env(), info, msg).unwrap();
        let res = query(deps.as_ref(), mock_env(), QueryMsg::GetActualSet { user: addr }).unwrap();
        let value: GetActualSetResponse = from_binary(&res).unwrap();
        assert_eq!("22222".to_string(), value.set.cards_id[0]);
        assert_eq!("".to_string(), value.set.cards_id[1]);

        // delete exist card from set 
        let addr =Addr::unchecked("anyone");
        let info = mock_info("anyone", &coins(2, "token"));
        let msg = ExecuteMsg::SetCardToSet { num_in_set: 0, card_id: "".to_string() };
        let _res = execute(deps.as_mut(), mock_env(), info, msg).unwrap();
        let res = query(deps.as_ref(), mock_env(), QueryMsg::GetActualSet { user: addr }).unwrap();
        let value: GetActualSetResponse = from_binary(&res).unwrap();
        assert_eq!("".to_string(), value.set.cards_id[0]);
        assert_eq!("".to_string(), value.set.cards_id[1]);

        
    }

    
    #[test]
    fn set_user_attribute() {
        let mut deps = mock_dependencies();

        let msg = InstantiateMsg {};
        let info = mock_info("creator", &coins(2, "token"));
        let _res = instantiate(deps.as_mut(), mock_env(), info, msg).unwrap();

        let user_param: UserParams=UserParams { ampers: 11, volts: 22 };
        let user_param_def: UserParams=UserParams { ampers: 0, volts: 0 };

        // beneficiary can release it
        let info = mock_info("anyone", &coins(2, "token"));
        let msg = ExecuteMsg::SetUserAttribute { num_in_set: 0, value: user_param };
        let _res = execute(deps.as_mut(), mock_env(), info, msg).unwrap();

        // create new card set
        let addr =Addr::unchecked("anyone");
        let res = query(deps.as_ref(), mock_env(), QueryMsg::GetActualSet { user: addr }).unwrap();
        let value: GetActualSetResponse = from_binary(&res).unwrap();
        assert_eq!(user_param, value.set.user_influence[0]);
        assert_eq!(user_param_def, value.set.user_influence[1]);        
    }
}
