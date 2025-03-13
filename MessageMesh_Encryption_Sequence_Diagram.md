```mermaid
sequenceDiagram
    title MessageMesh Encryption/Decryption Sequence

    %% Participants
    participant Client
    participant Network
    participant GetSymmetricKey
    participant Blockchain
    participant KeyPair
    participant SendFirstMessage
    participant EncryptDecrypt
    participant Base64

    %% EncryptMessage Function
    rect rgb(240, 240, 255)
    Note over Client,Base64: EncryptMessage Function
    
    Client->>Network: EncryptMessage(message, receiver)
    Network->>Network: Get sender ID and sort peerIDs
    Network->>GetSymmetricKey: GetSymmetricKey(peerIDs)
    
    alt Symmetric key found
        GetSymmetricKey-->>Network: Return symmetricKey
    else Symmetric key not found
        GetSymmetricKey-->>Network: Return nil
        Network->>Blockchain: CheckPeerFirstMessage(peerIDs)
        Network->>KeyPair: ReadKeyPair()
        KeyPair-->>Network: Return keyPair
        
        alt First message found in blockchain
            Blockchain-->>Network: Return firstMessage
            Network->>KeyPair: DecryptWithPrivateKey(firstMessage.GetSymetricKey(sender))
            KeyPair-->>Network: Return decrypted symmetricKey
        else First message not found
            Blockchain-->>Network: Return nil
            Network->>SendFirstMessage: SendFirstMessage(peerIDs, receiver)
            
            SendFirstMessage->>SendFirstMessage: Check if receiver is online
            SendFirstMessage->>SendFirstMessage: GenerateSymmetricKey(32)
            SendFirstMessage->>SendFirstMessage: SaveSymmetricKey(symmetricKey, peerIDs)
            SendFirstMessage->>SendFirstMessage: Encrypt symmetricKey for both peers
            SendFirstMessage->>SendFirstMessage: Create FirstMessage object
            SendFirstMessage->>SendFirstMessage: Send FirstMessage via PubSub
            
            SendFirstMessage-->>Network: Return firstMessage
            Network->>KeyPair: DecryptWithPrivateKey(firstMessage.GetSymetricKey(sender))
            KeyPair-->>Network: Return decrypted symmetricKey
        end
    end
    
    Network->>EncryptDecrypt: EncryptWithSymmetricKey(message, symmetricKey)
    EncryptDecrypt-->>Network: Return encryptedMessage
    Network->>Base64: base64.StdEncoding.EncodeToString(encryptedMessage)
    Base64-->>Network: Return base64Message
    Network-->>Client: Return base64Message
    end

    %% DecryptMessage Function
    rect rgb(255, 240, 240)
    Note over Client,Base64: DecryptMessage Function
    
    Client->>Network: DecryptMessage(message, peerIDs)
    Network->>Network: Sort peerIDs
    Network->>GetSymmetricKey: GetSymmetricKey(peerIDs)
    
    alt Symmetric key found
        GetSymmetricKey-->>Network: Return symmetricKey
    else Symmetric key not found
        GetSymmetricKey-->>Network: Return nil
        Network->>Blockchain: CheckPeerFirstMessage(peerIDs)
        
        alt First message found in blockchain
            Blockchain-->>Network: Return firstMessage
            Network->>KeyPair: ReadKeyPair()
            KeyPair-->>Network: Return keyPair
            Network->>KeyPair: DecryptWithPrivateKey(firstMessage.GetSymetricKey(selfID))
            KeyPair-->>Network: Return decrypted symmetricKey
            Network->>Network: SaveSymmetricKey(symmetricKey, peerIDs)
        else First message not found
            Blockchain-->>Network: Return nil
            Network-->>Client: Return error "first message not found"
        end
    end
    
    Network->>Base64: base64.StdEncoding.DecodeString(message)
    Base64-->>Network: Return encryptedBytes
    Network->>EncryptDecrypt: DecryptWithSymmetricKey(encryptedBytes, symmetricKey)
    EncryptDecrypt-->>Network: Return decryptedMessage
    Network-->>Client: Return decryptedMessage as string
    end 