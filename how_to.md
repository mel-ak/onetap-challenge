# How to Use the Bill Aggregation Service

## Overview
This service allows you to aggregate and manage your utility bills from multiple providers in one place. You can link accounts from different utility providers (electricity, water, internet, etc.) and view all your bills in a single dashboard.

## Prerequisites
- Docker and Docker Compose installed on your system
- Go 1.x (if running locally)
- Access to utility provider accounts you want to link

## Getting Started

### 1. Setup
1. Clone the repository
2. Run the following command to start the service:
```bash
make up
```
This will start all necessary services including the database and API server.

### 2. Authentication
Before using the service, you need to authenticate:
1. Register a new account using the registration endpoint
2. Login to receive your authentication token
3. Include this token in the `Authorization` header for all subsequent requests

## Using the Service

### Linking Utility Accounts
To link a utility account:
1. Navigate to the "Link Account" section
2. Select your utility provider
3. Enter your provider account credentials
4. Submit the request

Example API call:
```bash
curl -X POST /accounts/link \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "electricity_provider",
    "credentials": {
      "username": "your_username",
      "password": "your_password"
    }
  }'
```

### Viewing Bills
You can view your bills in several ways:

1. **View All Bills**
   - Access the main dashboard to see all your bills
   - Bills are automatically sorted by due date
   - Total amount due is calculated and displayed

2. **View Bills by Provider**
   - Filter bills by specific utility provider
   - See detailed breakdown of charges

Example API call:
```bash
curl -X GET /bills \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Refreshing Bill Data
To ensure you have the latest bill information:
1. Use the refresh button on the dashboard
2. Or make an API call to refresh bills

Example API call:
```bash
curl -X POST /bills/refresh \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Managing Linked Accounts
You can manage your linked accounts:
1. View all linked accounts
2. Remove accounts you no longer want to track
3. Update account credentials if needed

Example API call to remove an account:
```bash
curl -X DELETE /accounts/{account_id} \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Best Practices
1. **Regular Updates**: Refresh your bills regularly to ensure you have the latest information
2. **Account Security**: Keep your provider credentials secure and update them if compromised
3. **Bill Monitoring**: Set up notifications for upcoming due dates
4. **Data Accuracy**: Verify bill amounts against provider statements periodically

## Troubleshooting
If you encounter issues:
1. Check your internet connection
2. Verify your authentication token is valid
3. Ensure your provider credentials are correct
4. Check if the service is running properly
5. Contact support if issues persist

## API Documentation
For detailed API documentation, visit the Swagger UI at:
```
http://localhost:8080/swagger/index.html
```

## Support
If you need help:
1. Check the documentation
2. Review the troubleshooting guide
3. Contact support with specific error messages
4. Check the service status page for any ongoing issues
