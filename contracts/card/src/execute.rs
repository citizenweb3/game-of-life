use std::clone;

use cosmwasm_std::{Addr, CosmosMsg, Deps, DepsMut, Env, MessageInfo};
use cw721::Cw721Query;
use cw721_base::MintMsg;

use crate::error::ContractError;
use crate::helpers::generate_random_by_level;
use crate::state::{ Card, CONFIG, CardContract, Extension };
use cyber_std::CyberMsgWrapper;
type Response = cosmwasm_std::Response<CyberMsgWrapper>;


pub fn execute_execute(
    deps: DepsMut,
    _env: Env,
    info: MessageInfo,
    msgs: Vec<CosmosMsg<CyberMsgWrapper>>,
) -> Result<Response, ContractError> {
    let mut res = Response::new().add_attribute("action", "execute");

    let config = CONFIG.load(deps.storage)?;
    let owner = config.owner;
    if info.sender != owner {
        return Err(ContractError::Unauthorized {});
    }

    res = res.add_messages(msgs);

    Ok(res)
}
pub fn execute_mint(
    _deps: DepsMut,
    _env: Env,
    _info: MessageInfo,
    _mint_msg: MintMsg<Extension>,
) -> Result<Response, ContractError> {
    Err(ContractError::DisabledFunctionality {})
}

pub fn create_card(deps: DepsMut, env: Env, info: MessageInfo) -> Result<Response, ContractError>  {
    let executor_addr: String = info.clone().sender.into();

    for _i in 0..5 {
        let card_param = generate_random_by_level(0 as u8);
        let avatar_url: String = get_avatar_url();
        let card = Card{ param: card_param, avatar: avatar_url };
        let token_id = card.get_key(); 
        if is_uniq_card(deps.as_ref(), env.clone(), token_id) {
            return save_card(deps, env, info, executor_addr.clone(), card);
        }
    }
    
    return Err(ContractError::CardNotCreated {  });
}


fn get_avatar_url() -> String {
    return "AVATAR_URL is not realaze now".to_string();
}

fn is_uniq_card(deps: Deps, env: Env, token_id: String) -> bool{
    let cw721_contract = CardContract::default();

    let res = cw721_contract.owner_of(deps, env, token_id.clone(), true);
    
    return res.is_err();
}

fn save_card(deps: DepsMut, env: Env, info: MessageInfo, owner: String, card: Card) -> Result<Response, ContractError>  {
    let token_id: String =  card.get_key();
    let mint_msg = MintMsg {
        token_id: token_id,
        owner: owner,
        token_uri: None,
        extension: card,
    };
    // mint card by basic way in NFT
    let res = CardContract::default().mint(deps, env, info, mint_msg);
    if res.is_err() {
        return Err(ContractError::CardNotCreated{ });
    }

    return Ok(res?);
}

