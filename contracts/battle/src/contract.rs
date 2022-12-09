#[cfg(not(feature = "library"))]
use cosmwasm_std::entry_point;
use cosmwasm_std::{to_binary, Addr, Binary, Deps, DepsMut, Env, MessageInfo, Response, StdResult};

use cw2::set_contract_version;

use crate::error::ContractError;
use crate::msg::{ExecuteMsg, IsOpenToBattleResponse, InstantiateMsg, QueryMsg};
use crate::state::READY_TO_BATTLE;

// version info for migration info
const CONTRACT_NAME: &str = "crates.io:battle";
const CONTRACT_VERSION: &str = env!("CARGO_PKG_VERSION");

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn instantiate(
    deps: DepsMut,
    _env: Env,
    info: MessageInfo,
    msg: InstantiateMsg,
) -> Result<Response, ContractError> {
    set_contract_version(deps.storage, CONTRACT_NAME, CONTRACT_VERSION)?;

    Ok(Response::new()
        .add_attribute("method", "instantiate")
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
        ExecuteMsg::SetReadyToBattleStatus {ready} => execute::set_ready_to_battle_status(deps, info, ready),
        ExecuteMsg::Battle { rival } => execute::battle(deps, info, rival),
    }
}

pub mod execute {

    use super::*;

    pub fn set_ready_to_battle_status(deps: DepsMut, info: MessageInfo, ready: bool) -> Result<Response, ContractError> {
        READY_TO_BATTLE.update(deps.storage, info.sender.clone(), |_| -> Result<_, ContractError> {
            Ok(ready)
        })?;

        Ok(Response::new().
        add_attribute("val", ready.to_string()).
        add_attribute("action", "set_ready_status"))
    }

    pub fn battle(deps: DepsMut, info: MessageInfo, rival: Addr) -> Result<Response, ContractError> {
        if !READY_TO_BATTLE.load(deps.storage, rival.clone()).unwrap_or_default() {
            return Err(ContractError::RivalIsNotReady{});
        }
        // Get Set executor
        // Get user param (amper, volts...) executor
        // Modify card set executor

        
        // Get Set rival
        // Get user param (amper, volts...) rival
        // Modify card set rival

        // do battle(modified card sets executor and rival)

        Ok(Response::new().
        add_attribute("action", "battle"))
    }
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn query(deps: Deps, _env: Env, msg: QueryMsg) -> StdResult<Binary> {
    match msg {
        QueryMsg::IsOpenToBattle {user_id} => to_binary(&query::is_open_to_battle(deps, user_id)?),
    }
}

pub mod query {
    use super::*;

    pub fn is_open_to_battle(deps: Deps, user_id: Addr) -> StdResult<IsOpenToBattleResponse> {
        let state = READY_TO_BATTLE.load(deps.storage, user_id.clone()).unwrap_or_default();
        Ok(IsOpenToBattleResponse { ready: state})
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

        let msg = InstantiateMsg {};
        let info = mock_info("creator", &coins(1000, "earth"));

        // we can just call .unwrap() to assert this was a success
        let res = instantiate(deps.as_mut(), mock_env(), info, msg).unwrap();
        assert_eq!(0, res.messages.len());

        // it worked, let's query the state
        let res = query(deps.as_ref(), mock_env(), QueryMsg::IsOpenToBattle { user_id: Addr::unchecked("creator") }).unwrap();
        let value: IsOpenToBattleResponse = from_binary(&res).unwrap();
        assert_eq!(false, value.ready);
    }

    #[test]
    fn set_ready_to_battle_status() {
        let mut deps = mock_dependencies();

        let msg = InstantiateMsg {};
        let info = mock_info("creator", &coins(2, "token"));
        let _res = instantiate(deps.as_mut(), mock_env(), info, msg).unwrap();

        
        let user_test_str = "anymore";
        let user_test = Addr::unchecked(user_test_str);
        let info = mock_info(user_test_str, &coins(2, "token"));
        // check default status (must be false)
        let res = query(deps.as_ref(), mock_env(), QueryMsg::IsOpenToBattle { user_id: user_test.clone()}).unwrap();
        let value: IsOpenToBattleResponse = from_binary(&res).unwrap();
        assert_eq!(false, value.ready);

        // set status open
        let msg = ExecuteMsg::SetReadyToBattleStatus { ready: true};
        let _res = execute(deps.as_mut(), mock_env(), info.clone(), msg).unwrap();

        // check status open
        let res = query(deps.as_ref(), mock_env(), QueryMsg::IsOpenToBattle { user_id: user_test.clone()}).unwrap();
        let value: IsOpenToBattleResponse = from_binary(&res).unwrap();
        assert_eq!(true, value.ready);
        
        // update status to close
        let msg = ExecuteMsg::SetReadyToBattleStatus { ready: false};
        let _res = execute(deps.as_mut(), mock_env(), info.clone(), msg).unwrap();

        // check status close
        let res = query(deps.as_ref(), mock_env(), QueryMsg::IsOpenToBattle { user_id: user_test.clone()}).unwrap();
        let value: IsOpenToBattleResponse = from_binary(&res).unwrap();
        assert_eq!(false, value.ready);
    }

}
