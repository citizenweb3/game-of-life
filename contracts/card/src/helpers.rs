use std::u8;

use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

use cosmwasm_std::{
    to_binary, Addr, CosmosMsg, CustomQuery, Querier, QuerierWrapper, StdResult, WasmMsg, WasmQuery,
};

use crate::msg::{ExecuteMsg, QueryMsg};
use crate::state::{CardParams};
use rand::prelude::*;

/// CwTemplateContract is a wrapper around Addr that provides a lot of helpers
/// for working with this.
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct CwTemplateContract(pub Addr);

impl CwTemplateContract {
    pub fn addr(&self) -> Addr {
        self.0.clone()
    }

    pub fn call<T: Into<ExecuteMsg>>(&self, msg: T) -> StdResult<CosmosMsg> {
        let msg = to_binary(&msg.into())?;
        Ok(WasmMsg::Execute {
            contract_addr: self.addr().into(),
            msg,
            funds: vec![],
        }
        .into())
    }
}

pub fn generate_random_by_level(level: u8) -> CardParams {
    let mult = get_level_multiplier(level);
    let card_params = CardParams {
      hp: (get_random_between(100, 100) as f32 * mult) as u8,
      level: level,
      accuracy: (get_random_between(10, 20) as f32 * mult) as u8 ,
      damage: (get_random_between(20, 30) as f32 * mult) as u8,
      deffence: (get_random_between(0, 20) as f32 * mult) as u8,
    };
    return card_params;
}


fn get_level_multiplier(level: u8) -> f32 {
    return 1.0 + (0.1 * level as f32);
}

fn get_random_between(min: u8, range: u8) -> u8 {
    let res = get_random_u8()%range;
    if u8::MAX - min < res {
        return  u8::MAX;
    }
    return min + res;
}
  
fn get_random_u8() -> u8 {
    let mut rng = rand::thread_rng();
    let n: u8 = rng.gen();
    n
}


#[cfg(test)]
mod tests {
    use super::*;
    use std::collections::HashMap;
    
    #[test]
    fn test_get_random_u8() {
        let r1 = get_random_u8();
        let r2 = get_random_u8();
        assert_ne!(r1, r2);
    }

    
    #[test]
    fn test_get_random_between() {
        let r0_100 = get_random_between(0, 100);
        assert!(r0_100 <= 100, "rand(0, 100) = {}", r0_100);
        
        // check out of range
        let _ = get_random_between(u8::MAX, u8::MAX);
        
        let min = get_random_u8();
        let range = get_random_u8();
        let r2 = get_random_between(min, range);
        
        assert!(r2 >= min);
        assert!(r2-min <= range);
    }

    #[test]
    fn test_get_level_multiplier() {
        // arrange
        let mut map = HashMap::new();
        map.insert(0,1.0 );
        map.insert(1,1.1 );
        map.insert(2,1.2 );
        map.insert(10,2.0 );
        map.insert(100,11.0 );

        for (level, expected_mult) in map {
            // action 
            let actual_mult = get_level_multiplier(level); 
            // assert 
            assert_eq!(expected_mult, actual_mult);  
        }
    }

}