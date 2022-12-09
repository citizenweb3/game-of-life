use cosmwasm_schema::{cw_serde, QueryResponses};
use cosmwasm_std::Addr;

use crate::state::Card;


#[cw_serde]
pub struct InstantiateMsg {
    pub count: i32,
}

#[cw_serde]
pub enum ExecuteMsg {
    Mint {},
    Transfer { card_id: String, to: Addr },
    
    // todo realize in future
    // Burn { cardId: String},
    // Freeze { cardId1: String, cardId2: String},
    // UnFreeze { cardId: String},
}

#[cw_serde]
#[derive(QueryResponses)]
pub enum QueryMsg {
    #[returns(GetCardResponse)]
    GetCard{card_id: String},
    
    #[returns(GetCardOwnerResponse)]
    GetCardOwner {card_id: String},
    
    #[returns(GetOwnersCardsResponse)]
    GetOwnersCards {user_id: Addr},
    
    #[returns(IsOwnerResponse)]
    IsOwner {card_id: String, user_id: Addr},
}

// We define a custom struct for each query response
#[cw_serde]
pub struct GetCardResponse {
    pub card: Card,
}

#[cw_serde]
pub struct GetCardOwnerResponse {
    pub user_id: Addr,
}

#[cw_serde]
pub struct GetOwnersCardsResponse {
    pub card_ids: Vec<String>, // todo change to &[String]
}

#[cw_serde]
pub struct IsOwnerResponse {
    pub is_owner: bool,
}

