    ```mermaid
    sequenceDiagram
        participant Client
        participant Server
        Client->>Server: Request data
        activate Server
        Server-->>Client: Send data
        deactivate Server
