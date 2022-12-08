use cosmwasm_std::StdError;
use thiserror::Error;

#[derive(Error, Debug)]
pub enum ContractError {
    #[error("{0}")]
    Std(#[from] StdError),

    #[error("Unauthorized")]
    Unauthorized {},
    // Add any other custom errors you like here.
    // Look at https://docs.rs/thiserror/1.0.21/thiserror/ for details.

    
    #[error("Num more than count count in set")]
    TooMuchCardNum {},
    
    #[error("Total sum one of the param more than 100")]
    TooMuchInfluense {},
}
