use cosmwasm_schema::{cw_serde, QueryResponses};
use cosmwasm_std::Addr;

#[cw_serde]
pub struct InstantiateMsg {
}

#[cw_serde]
pub enum ExecuteMsg {
    SetReadyToBattleStatus{ready: bool},
    Battle{rival: Addr},
}

#[cw_serde]
#[derive(QueryResponses)]
pub enum QueryMsg {
    #[returns(IsOpenToBattleResponse)]
    IsOpenToBattle{user_id: Addr},

    // #[returns(GetModifedCardsResponse)]
    // GetModifedCards{},
}

// Custom struct for each query response
#[cw_serde]
pub struct IsOpenToBattleResponse {
    pub ready: bool,
}

// #[cw_serde]
// pub struct GetModifedCardsResponse {

// }