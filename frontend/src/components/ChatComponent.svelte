<script lang="ts">
  import { Badge, Indicator, Textarea, ToolbarButton } from 'flowbite-svelte';
  import { Navbar, NavBrand, NavLi, NavUl, NavHamburger } from 'flowbite-svelte';
  import { PaperPlaneOutline } from 'flowbite-svelte-icons';
  import { SendMessage } from '../../wailsjs/go/main/App.js';
  import * as Wails from '../../wailsjs/runtime/runtime.js';

  let { userPeerID = $bindable() } = $props();
  let message = $state('');
  // let chatService = new backend.ChatService();
  Wails.EventsOn('getMessage', (newMessage) => {
    messageList.push(newMessage);
    console.log(messageList);
  });
  let messageList = $state([{}]);

  function sendMessage(): void {
    SendMessage(message);
    message = '';
  }
</script>

<div class="flex flex-col h-full flex-auto">
  <div class="flex-none">
    <Navbar>
      <NavBrand href="#">
        <span class="self-center whitespace-nowrap text-xl font-semibold dark:text-white">Flowbite</span>
      </NavBrand>
      <NavHamburger />
      <NavUl>
        <Badge color="red" rounded class="px-2.5 py-0.5">
          <Indicator color="red" size="xs" class="me-1" />Unavailable
        </Badge>
      </NavUl>
    </Navbar>
  </div>
  <div class="flex-auto flex flex-col items-end h-full overflow-y-auto">
    <!-- Check if message is from self or other -->
    {#each messageList as message}
      {#if message.sender === userPeerID}
      <div class="flex w-full justify-end p-3">
        <div class="flex flex-col w-full max-w-[320px] leading-1.5 p-4 text-white bg-primary-700 dark:bg-primary-800 rounded-l-xl rounded-br-xl">
          <div class="flex items-center space-x-2 rtl:space-x-reverse">
            <span class="text-sm font-semibold text-white">{message.sender}</span>
            <span class="text-sm font-normal text-gray-300">11:46</span>
          </div>
          <p class="text-sm font-normal py-2.5 text-white">{message.message}</p>
          <span class="text-sm font-normal text-end text-gray-300">Delivered</span>
        </div>
      </div>
      {:else}
      <div class="flex w-full p-3 mt-auto">
        <div class="flex flex-col w-full max-w-[320px] leading-1.5 p-4 border-gray-200 bg-gray-100 rounded-e-xl rounded-es-xl dark:bg-gray-700">
          <div class="flex items-center space-x-2 rtl:space-x-reverse">
            <span class="text-sm font-semibold text-gray-900 dark:text-white">{message.sender}</span>
            <span class="text-sm font-normal text-gray-500 dark:text-gray-400">11:46</span>
          </div>
          <p class="text-sm font-normal py-2.5 text-gray-900 dark:text-white">{message.message}</p>
          <span class="text-sm font-normal text-gray-500 dark:text-gray-400">Delivered</span>
        </div>
      </div>
      {/if}
    {/each}
  </div>
  <div class="flex-none">
    <label for="chat" class="sr-only">Your message</label>
    <div class="flex items-center px-3 py-2 rounded-none bg-gray-50 dark:bg-gray-700">
      <Textarea bind:value={message} id="chat" class="mx-4 bg-white dark:bg-gray-800 h-10 min-h-10 max-h-20" placeholder="Your message..." />
      <ToolbarButton onclick={sendMessage} color="blue" class="rounded-full text-primary-600 dark:text-primary-500">
        <PaperPlaneOutline class="w-6 h-6 rotate-45" />
        <span class="sr-only">Send message</span>
      </ToolbarButton>
    </div>
  </div>
</div>

<style lang="less">
</style>
