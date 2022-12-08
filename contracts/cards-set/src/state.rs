use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

use cosmwasm_std::Addr;
use cw_storage_plus::Map;

pub const COUNT_CARDS_IN_SET: usize = 5;

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct State {
    pub count: i32,
    pub owner: Addr,
}

#[derive(Serialize, Deserialize, Default, Copy, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct UserParams {
    pub ampers: u8,
    pub volts: u8,
    // todo Add another params
}

#[derive(Serialize, Deserialize, Default, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct CardsSet {
    pub cards_id: [String; 5],
    pub user_influence: [UserParams; 5]
}


pub const CARDS_SETS: Map<Addr, CardsSet> = Map::new("cards_sets");
