<script lang="ts">
  import './app.css';
  import NavigationRailComponent from './components/NavigationRailComponent.svelte';
  import ChatListComponent from './components/ChatListComponent.svelte';
  import ChatComponent from './components/ChatComponent.svelte';
  import * as Wails from '../wailsjs/runtime/runtime.js';
  import { models } from '../wailsjs/go/models.js';
  import { GetMessagesFromPeer, GetDecryptedMessage } from '../wailsjs/go/main/App.js';
    import { Input, Spinner, ToolbarButton } from 'flowbite-svelte';
    import { PaperPlaneOutline } from 'flowbite-svelte-icons';

  let selectedPeer = $state('');
  let userPeerID = $state('');
  let online = $state(false);
  let peerList = $state<string[]>([]); // All peers
  let onlinePeerList = $state<string[]>([]); // Peers that are online
  let messages = $state<models.Message[]>([]);
  let accounts = $state<models.Account[]>([]);
  let messageMap = $state(new Map<string, models.Block[]>());
  let accountMap = $state(new Map<string, models.Account>());
  let topic = $state('');
  let topicChanged = $state(false);

  Wails.EventsOn("getPeerList", (data: string[]) => {
    onlinePeerList = data ?? [];
    peerList = [...new Set([...peerList, ...onlinePeerList])];
  });
  Wails.EventsOn("getUserPeerID", (data: string) => {
    userPeerID = data;
  });

  Wails.EventsOn("getAccounts", (data: models.Account[]) => {
    accounts = data;
  });

  Wails.EventsOn("getConnected", (data: boolean) => {
    online = data;
  });

  // Load initial blockchain data
  Wails.EventsOn("getBlockchain", (blocks: models.Block[]) => {
    console.log("getBlockchain", blocks);
    blocks.forEach(block => {
      if (block.BlockType === "message") {
        const message: models.Message = block.Data;
        const key = getMessageKey(message.sender, message.receiver);
        if (!messageMap.has(key)) {
          messageMap.set(key, []);
        } if (!messageMap.get(key)?.some(m => m.Hash === block.Hash)) {
          messageMap.get(key)?.push(block);
        }
        if (selectedPeer && (message.sender === selectedPeer || message.receiver === selectedPeer)) {
          getMessagesForPeer([selectedPeer, userPeerID]).then(msgs => {
            messages = msgs;
          });
        }
      } else if (block.BlockType === "account") {
        const account: models.Account = block.Data;
        accountMap.set(account.publicKey, account);
      }
    });
  });

  Wails.EventsOn("getBlock", (block: models.Block) => {
    if (block.BlockType === "message") {
      const message: models.Message = block.Data;
      const key = getMessageKey(message.sender, message.receiver);
      if (!messageMap.has(key)) {
        messageMap.set(key, []);
      } if (!messageMap.get(key)?.some(m => m.Hash === block.Hash)) {
        messageMap.get(key)?.push(block);
      }
    }
    if (block.BlockType === "account") {
      const account: models.Account = block.Data;
      accountMap.set(account.publicKey, account);
    }
  });

  function getMessageKey(sender: string, receiver: string): string {
    // Create consistent key regardless of sender/receiver order
    return [sender, receiver].sort().join(':');
  }

  async function getMessagesForPeer(peerIDs: string[]): Promise<models.Message[]> {
    const messages: models.Message[] = [];
    for (const [key, msgs] of messageMap.entries()) {
      if (getMessageKey(peerIDs[0], peerIDs[1]) === key) {
        for (const msg of msgs) {
          const decryptedMessage = await GetDecryptedMessage(msg.Data.message, peerIDs);
          messages.push({
            ...msg.Data,
            message: decryptedMessage
          } as models.Message);
        }
      }
    }
    return messages.sort((a, b) => 
      new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime()
    );
  }

  // Bind these to child components
  $effect(() => {
    if (selectedPeer) {
      messages = [];
      getMessagesForPeer([selectedPeer, userPeerID]).then(msgs => {
        messages = msgs;
      });
    }
  });

</script>

<main>
  <div class="flex w-screen h-screen bg-primary-50 dark:bg-gray-900">
    <!-- If online is false, show a loading screen -->
    {#if !online && !topicChanged}
      <div class="flex flex-col items-center justify-center h-full w-full">
        <div class="text-2xl font-bold">Welcome to MessageMesh</div>
        <div class="text-sm text-gray-500">Please enter a topic to start chatting</div>
        <div class="flex flex-row">
          <Input type="text" bind:value={topic} class="w-full" />
          <ToolbarButton class="bg-blue-500 text-white p-2 rounded-md" on:click={() => {
            if (topic !== "") {
              topicChanged = true;
              Wails.EventsEmit("joinTopic", topic);
            }
          }}>
            <PaperPlaneOutline class="w-6 h-6 rotate-45" />
            <span class="sr-only">Join Topic</span>
          </ToolbarButton>
        </div>
      </div>
    {/if}
    {#if !online && topicChanged}
      <div class="flex flex-col items-center justify-center h-full w-full">
        <Spinner size={8} />
        <div class="text-2xl font-bold">Connecting to network...</div>
      </div>
    {/if}
    {#if online && topicChanged}
    <NavigationRailComponent bind:onlinePeerList bind:online bind:userPeerID bind:topic></NavigationRailComponent>
      <div class="flex flex-row w-full">
        <ChatListComponent bind:selectedPeer bind:peerList bind:onlinePeerList></ChatListComponent>
        <ChatComponent bind:userPeerID bind:selectedPeer bind:messages></ChatComponent>
      </div>
    {/if}
  </div>
</main>

<style>
</style>
