use cosmwasm_std::Addr;
use cw_storage_plus::Map;

pub const READY_TO_BATTLE: Map<Addr, bool> = Map::new("ready_to_battle");
