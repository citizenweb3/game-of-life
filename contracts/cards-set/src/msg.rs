use cosmwasm_schema::{cw_serde, QueryResponses};
use crate::state::{UserParams, CardsSet};

use cosmwasm_std::Addr;

#[cw_serde]
pub struct InstantiateMsg {
}

#[cw_serde]
pub enum ExecuteMsg {
    SetCardToSet {num_in_set: usize, card_id: String},
    SetUserAttribute {num_in_set: usize, value: UserParams},
}

#[cw_serde]
#[derive(QueryResponses)]
pub enum QueryMsg {
    
    #[returns(GetActualSetResponse)]
	GetActualSet{user: Addr},

    
    // use GetActualSetResponse. User Attribute is included in set.
    // #[returns(GetUserAttributeResponse)]
	// GetUserAttribute{},

    
    // #[returns(GetActualSetWithAttributeResponse)]
	// GetActualSetWithAttribute{},

}

// We define a custom struct for each query response
#[cw_serde]
pub struct GetActualSetResponse {
    pub set: CardsSet,
}

