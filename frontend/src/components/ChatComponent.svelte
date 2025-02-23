<script lang="ts">
  import { Badge, Indicator, Textarea, ToolbarButton } from 'flowbite-svelte';
  import { Navbar, NavBrand, NavLi, NavUl, NavHamburger } from 'flowbite-svelte';
  import { PaperPlaneOutline } from 'flowbite-svelte-icons';
  import * as Wails from '../../wailsjs/runtime/runtime.js';
  import { SendMessage, GetMessagesFromPeer } from '../../wailsjs/go/main/App.js';
  import { models } from '../../wailsjs/go/models.js';
  let { userPeerID = $bindable(), selectedPeer = $bindable(), messages = $bindable() } = $props();
  
  let message = $state('');
  let messagesContainer: HTMLDivElement;
  
  function scrollToBottom(): void {
    if (messagesContainer) {
      messagesContainer.scrollTop = messagesContainer.scrollHeight;
    }
  }
  
  $effect(() => {
    if (messages) {
      setTimeout(scrollToBottom, 0);
    }
  });

  function sendMessage(): void {
    if (!selectedPeer) return; // Don't send if no peer is selected
    SendMessage(message, selectedPeer);
    message = '';
  }
</script>

<div class="flex flex-col h-screen flex-auto">
  <!-- Fixed navbar at top -->
  <div id="navbar" class="flex-none">
    <Navbar>
      <NavBrand href="#">
        <span class="self-center whitespace-nowrap text-xl font-semibold text-ellipsis dark:text-white">
          {selectedPeer || 'Select a chat'}
        </span>
      </NavBrand>
      <!-- <NavHamburger />
      <NavUl>
        <Badge color="red" rounded class="px-2.5 py-0.5">
          <Indicator color="red" size="xs" class="me-1" />Unavailable
        </Badge>
      </NavUl> -->
    </Navbar>
  </div>

  <!-- Scrollable messages area -->
  <div id="messages" class="flex-1 overflow-hidden">
    <div bind:this={messagesContainer} class="h-full overflow-y-auto">
      <!-- Check if message is from self or other -->
      {#each messages as message}
        {#if message.sender === userPeerID }
        <div class="flex w-full justify-end p-3">
          <div class="flex flex-col w-full max-w-[320px] leading-1.5 p-4 text-white bg-primary-700 dark:bg-primary-800 rounded-l-xl rounded-br-xl">
            <span class="text-sm font-semibold text-white flex-initial text-ellipsis">{message.sender}</span>
            <p class="text-sm font-normal py-2.5 text-white">{message.message}</p>
            <span class="text-sm font-normal text-end text-gray-300">{new Date(message.timestamp).toLocaleTimeString()}</span>
          </div>
        </div>
        {:else}
        <div class="flex w-full p-3 mt-auto">
          <div class="flex flex-col w-full max-w-[320px] leading-1.5 p-4 border-gray-200 bg-gray-100 rounded-e-xl rounded-es-xl dark:bg-gray-700">
              <span class="text-sm font-semibold text-gray-900 flex-initial text-ellipsis dark:text-white">{message.sender}</span>
            <p class="text-sm font-normal py-2.5 text-gray-900 dark:text-white">{message.message}</p>
            <span class="text-sm font-normal text-gray-500 dark:text-gray-400">{new Date(message.timestamp).toLocaleTimeString()}</span>
          </div>
        </div>
        {/if}
      {/each}
    </div>
  </div>

  <!-- Fixed message input at bottom -->
  <div id="message-input" class="flex-none">
    <label for="chat" class="sr-only">Your message</label>
    <div class="flex items-center px-3 py-2 rounded-none bg-gray-50 dark:bg-gray-700">
      <Textarea 
        bind:value={message} 
        id="chat" 
        class="mx-4 bg-white dark:bg-gray-800 h-10 min-h-10 max-h-20" 
        placeholder="Your message..." 
      />
      <ToolbarButton 
        on:click={sendMessage} 
        color="blue" 
        class="rounded-full text-primary-600 dark:text-primary-500"
      >
        <PaperPlaneOutline class="w-6 h-6 rotate-45" />
        <span class="sr-only">Send message</span>
      </ToolbarButton>
    </div>
  </div>
</div>

<style lang="less">
  .text-ellipsis {
    text-overflow: ellipsis;
    overflow: hidden;
    white-space: nowrap;
  }
</style>
