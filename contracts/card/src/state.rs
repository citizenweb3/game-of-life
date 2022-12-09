use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

use cosmwasm_std::Addr;
use cw_storage_plus::Map;

#[derive(Serialize, Deserialize, Default, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct CardParams {
	pub hp:       u8,
	pub level:    u8,
	pub deffence: u8,
    pub damage: u8,
	pub accuracy: u8,
}
impl CardParams {
    pub fn empty() -> CardParams {
        return CardParams{
            hp:0, 
            level:0, 
            deffence:0, 
            damage:0, 
            accuracy:0,
        };
    }
    
}

#[derive(Serialize, Deserialize, Clone, Default, Debug, PartialEq, Eq, JsonSchema)]
pub struct Card {
    pub param: CardParams,
    pub avatar: String,
}


impl Card {
    pub fn empty() -> Card {
        return Card {
            param: CardParams::empty(),
            avatar: "".to_string(),
        };
    }
    pub fn get_key(&self) -> String {
        return self.param.hp.to_string();
    }
}

pub const CARDS: Map<String, Card> = Map::new("cards");
pub const OWNER_CARDS: Map<Addr, Vec<String>> = Map::new("owner_cards");
pub const CARD_OWNER: Map<String, Addr> = Map::new("card_owner");
