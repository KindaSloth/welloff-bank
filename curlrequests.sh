# script to test concurrent/parallel requests
function perform_transfer {
  curl --cookie "sessionId=01912a50-18ac-7859-944f-cb4745d34c7e; Max-Age=86400; Domain=localhost; Path=/; Secure; HttpOnly" \
       --request POST \
       --url http://localhost:5000/transaction/transfer \
       --header 'Content-Type: application/json' \
       --data '{
        "amount": "1.00",
        "from_account_id": "01911f41-f631-734b-a631-46aca9614536",
        "to_account_id": "01911fb7-727e-74e2-97d3-639eab3966ef"
       }'
}

export -f perform_transfer

num_requests=300

parallelism=20

seq $num_requests | xargs -n1 -P$parallelism -I{} bash -c 'perform_transfer'