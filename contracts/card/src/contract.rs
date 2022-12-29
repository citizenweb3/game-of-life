#[cfg(not(feature = "library"))]
use cosmwasm_std::entry_point;
use cosmwasm_std::{Binary, Deps, DepsMut, Env, MessageInfo, StdResult};
use cw2::set_contract_version;
use cyber_std::CyberMsgWrapper;


use crate::error::ContractError;
use crate::execute::{execute_mint, create_card};
use crate::msg::{ExecuteMsg, InstantiateMsg, QueryMsg};
use crate::state::{CardContract};

type Response = cosmwasm_std::Response<CyberMsgWrapper>;


// version info for migration info
const CONTRACT_NAME: &str = "game-of-life-card";
const CONTRACT_VERSION: &str = "0.0.1";

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn instantiate(
    mut deps: DepsMut,
    env: Env,
    info: MessageInfo,
    mut msg: InstantiateMsg,
) -> Result<Response, ContractError> {
    // override minter to contract itself
    msg.minter = env.clone().contract.address.into_string();
    let res = CardContract::default().instantiate(deps.branch(), env, info, msg.into())?;
    set_contract_version(deps.storage, CONTRACT_NAME, CONTRACT_VERSION)?;

    Ok(res)
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn execute(
    deps: DepsMut,
    env: Env,
    info: MessageInfo,
    msg: ExecuteMsg,
) -> Result<Response, ContractError> {
    match msg {
        // Freeze todo: realize this method
        // Unfreeze todo: realize this method
        ExecuteMsg::Mint(mint_msg) => execute_mint(deps, env, info, mint_msg),
        ExecuteMsg::CreateCard{} => create_card(deps, env, info),
        
        // CW721 methods
        _ => CardContract::default()
        .execute(deps, env, info, msg.into())
        .map_err(|err| err.into()),
    }
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn query(deps: Deps, env: Env, msg: QueryMsg) -> StdResult<Binary> {
    match msg {
        // QueryMsg::GetCard { card_id } realize serialized getter for nft???

        // CW721 methods
        _ => CardContract::default().query(deps, env, msg.into()),
    }
}

// pub mod query {
//     use super::*;
    
//     pub fn get_card(deps: Deps, card_id: String) -> StdResult<GetCardResponse> {
//         let card_info = CARDS.load(deps.storage, card_id).unwrap_or_default();

//         Ok(GetCardResponse { card: card_info})
//     }
// }
