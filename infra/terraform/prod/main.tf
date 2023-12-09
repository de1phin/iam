

module "iam" {
    source = "../modules/iam"

    folder_id = "b1gqe3skkuiko3bv671e"

    dns_domain = "iam.de1phin.ru"
    internal_dns_domain = "iam.internal"

    database = [
        {
        "dbname": "account-service"
        "user": "account-service-user"
        },
        {
        "dbname": "token-service"
        "user": "token-service-user"
        }
    ]

    dns_endpoints = [
        {
            "ip": "10.96.203.91"
            "hostname": "account"
            "public": false
        },
        {
            "ip": "158.160.133.43"
            "hostname": "account"
            "public": true
        },

        {
            "ip": "10.96.151.216"
            "hostname": "token"
            "public": false
        },
        {
            "ip": "158.160.131.19"
            "hostname": "token"
            "public": true
        }
    ]
}
