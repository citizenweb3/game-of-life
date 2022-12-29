
#[cfg(test)]
mod tests {
    use cosmwasm_std::testing::{mock_dependencies, mock_env, mock_info};
    use cosmwasm_std::from_binary;
    use crate::msg::{ExecuteMsg, InstantiateMsg, QueryMsg};
    use crate::contract::{execute, instantiate, query};

    use cw721::NumTokensResponse;

    #[test]
    fn proper_initialization() {
        let mut deps = mock_dependencies();
        let owner = "owner";
        let any_user = "bostrom19nk207agguzdvpj9nqsf4zrjw8mcuu9afun3fv";

        let msg = InstantiateMsg { 
            name: "GameOfLifeCard".to_string(),
            symbol:"GOLC".to_string(),
            minter: owner.to_string(),
            owner: owner.to_string(),
        };
        let info = mock_info(owner, &[]);

        // we can just call .unwrap() to assert this was a success
        let res = instantiate(deps.as_mut(), mock_env(), info, msg).unwrap();
        assert_eq!(0, res.messages.len());

        //
        let create_card_msg = ExecuteMsg::CreateCard { };
        let info = mock_info(&any_user, &[]);
        execute(deps.as_mut(), mock_env(), info, create_card_msg).unwrap();
        
        let get_count_tokens_msg = QueryMsg::NumTokens {  };
        let resp_count_token = query(deps.as_ref(), mock_env(), get_count_tokens_msg);

        assert!(resp_count_token.is_ok());
        let count = resp_count_token.unwrap();
        let res: NumTokensResponse = from_binary(&count).unwrap();
        let expected_num_token = NumTokensResponse{ count: 1 };
        assert_eq!(res, expected_num_token);
    }
}