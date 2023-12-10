

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
        },
        {
        "dbname": "access-service"
        "user": "access-service-user"
        }
    ]

    dns_endpoints = [
        {
            "ip": ""
            "hostname": "token"
            "public": true
        }
    ]
}
