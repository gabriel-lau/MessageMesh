<script lang="ts">
  import './app.css';
  import NavigationRailComponent from './components/NavigationRailComponent.svelte';
  import ChatListComponent from './components/ChatListComponent.svelte';
  import ChatComponent from './components/ChatComponent.svelte';
  import * as Wails from '../wailsjs/runtime/runtime.js';
  import { models } from '../wailsjs/go/models.js';
  import { GetMessagesFromPeer } from '../wailsjs/go/main/App.js';

  let selectedPeer = $state('');

  let userPeerID = $state('');
  Wails.EventsOn("getUserPeerID", (data: string) => {
    userPeerID = data;
  });
  let messages = $state<models.Message[]>([]);
  let ready = $state(false);
  Wails.EventsOn("ready", () => {
    ready = true;
  });
  let accounts = $state<models.Account[]>([]);
  Wails.EventsOn("getAccounts", (data: models.Account[]) => {
    accounts = data;
  });
  // Store all messages in a map with composite key "sender:receiver"
  let messageMap = $state(new Map<string, models.Message[]>());
  // Store all accounts in a map with peerID as key
  let accountMap = $state(new Map<string, models.Account>());

  // Load initial blockchain data
  Wails.EventsOn("getBlockchain", (blocks: models.Block[]) => {
    console.log("getBlockchain", blocks);
    blocks.forEach(block => {
      if (block.BlockType === "message") {
        const message = block.Data.Message;
        const key = getMessageKey(message.sender, message.receiver);
        if (!messageMap.has(key)) {
          messageMap.set(key, []);
        }
        messageMap.get(key)?.push(message);
      } else if (block.BlockType === "account") {
        const account = block.Data.Account;
        accountMap.set(account.publicKey, account);
      }
    });
  });

  // Listen for new messages
  Wails.EventsOn("getMessage", (message: models.Message) => {
    const key = getMessageKey(message.sender, message.receiver);
    if (!messageMap.has(key)) {
      messageMap.set(key, []);
    }
    messageMap.get(key)?.push(message);
    
    // Update messages if this message belongs to the selected peer
    if (selectedPeer && (message.sender === selectedPeer || message.receiver === selectedPeer)) {
      messages = getMessagesForPeer(selectedPeer);
    }
  });

  // Listen for new accounts
  Wails.EventsOn("getAccount", (account: models.Account) => {
    accountMap.set(account.publicKey, account);
  });

  function getMessageKey(sender: string, receiver: string): string {
    // Create consistent key regardless of sender/receiver order
    return [sender, receiver].sort().join(':');
  }

  function getMessagesForPeer(peerId: string): models.Message[] {
    const messages: models.Message[] = [];
    messageMap.forEach((msgs, key) => {
      if (key.includes(peerId)) {
        messages.push(...msgs);
      }
    });
    return messages.sort((a, b) => 
      new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime()
    );
  }

  // Bind these to child components
  messages = getMessagesForPeer(selectedPeer);
  $effect(() => {
    if (selectedPeer) {
      messages = getMessagesForPeer(selectedPeer);
    }
  });

</script>

<main>
  <div class="flex w-screen h-screen bg-primary-50 dark:bg-gray-900">
    <NavigationRailComponent></NavigationRailComponent>
    <div class="flex flex-row w-full">
      <ChatListComponent bind:userPeerID bind:selectedPeer bind:accounts></ChatListComponent>
      <ChatComponent bind:userPeerID bind:selectedPeer bind:messages></ChatComponent>
    </div>
  </div>
</main>

<style>
</style>
