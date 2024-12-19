<script lang="ts">
  import { Badge, Indicator, Textarea, ToolbarButton } from 'flowbite-svelte';
  import { Navbar, NavBrand, NavLi, NavUl, NavHamburger } from 'flowbite-svelte';
  import { PaperPlaneOutline } from 'flowbite-svelte-icons';

  let { isChatOpen = $bindable() } = $props();
  let message = $state('');

  function sendMessage() {
    console.log(message);
  }
</script>

<div class="flex flex-col h-full flex-auto">
  <div class="flex-none">
    <Navbar>
      <NavBrand href="/">
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
    <!-- Chat message 1 -->
    <div class="flex w-full p-3 mt-auto">
      <div class="flex flex-col w-full max-w-[320px] leading-1.5 p-4 border-gray-200 bg-gray-100 rounded-e-xl rounded-es-xl dark:bg-gray-700">
        <div class="flex items-center space-x-2 rtl:space-x-reverse">
          <span class="text-sm font-semibold text-gray-900 dark:text-white">Bonnie Green</span>
          <span class="text-sm font-normal text-gray-500 dark:text-gray-400">11:46</span>
        </div>
        <p class="text-sm font-normal py-2.5 text-gray-900 dark:text-white">That's awesome. I think our users will really appreciate the improvements.</p>
        <span class="text-sm font-normal text-gray-500 dark:text-gray-400">Delivered</span>
      </div>
    </div>
    <!-- Chat message 2 (Response) -->
    <div class="flex w-full justify-end p-3">
      <div class="flex flex-col w-full max-w-[320px] leading-1.5 p-4 border-gray-200 bg-gray-100 rounded-l-xl rounded-br-xl dark:bg-gray-700">
        <div class="flex items-center justify-end space-x-2 rtl:space-x-reverse">
          <span class="text-sm font-semibold text-gray-900 dark:text-white">Bonnie Green</span>
          <span class="text-sm font-normal text-gray-500 dark:text-gray-400">11:46</span>
        </div>
        <p class="text-sm font-normal py-2.5 text-gray-900 dark:text-white">That's awesome. I think our users will really appreciate the improvements.</p>
        <span class="text-sm font-normal text-gray-500 dark:text-gray-400">Delivered</span>
      </div>
    </div>
    <!-- Chat message 2 (Response) -->
    <div class="flex w-full justify-end p-3">
      <div class="flex flex-col w-full max-w-[320px] leading-1.5 p-4 border-gray-200 bg-gray-100 rounded-l-xl rounded-br-xl dark:bg-gray-700">
        <div class="flex items-center justify-end space-x-2 rtl:space-x-reverse">
          <span class="text-sm font-semibold text-gray-900 dark:text-white">Bonnie Green</span>
          <span class="text-sm font-normal text-gray-500 dark:text-gray-400">11:46</span>
        </div>
        <p class="text-sm font-normal py-2.5 text-gray-900 dark:text-white">That's awesome. I think our users will really appreciate the improvements.</p>
        <span class="text-sm font-normal text-gray-500 dark:text-gray-400">Delivered</span>
      </div>
    </div>
    <!-- Chat message 2 (Response) -->
    <div class="flex w-full justify-end p-3">
      <div class="flex flex-col w-full max-w-[320px] leading-1.5 p-4 border-gray-200 bg-gray-100 rounded-l-xl rounded-br-xl dark:bg-gray-700">
        <div class="flex items-center justify-end space-x-2 rtl:space-x-reverse">
          <span class="text-sm font-semibold text-gray-900 dark:text-white">Bonnie Green</span>
          <span class="text-sm font-normal text-gray-500 dark:text-gray-400">11:46</span>
        </div>
        <p class="text-sm font-normal py-2.5 text-gray-900 dark:text-white">That's awesome. I think our users will really appreciate the improvements.</p>
        <span class="text-sm font-normal text-gray-500 dark:text-gray-400">Delivered</span>
      </div>
    </div>
    <!-- Chat message 2 (Response) -->
    <div class="flex w-full justify-end p-3">
      <div class="flex flex-col w-full max-w-[320px] leading-1.5 p-4 border-gray-200 bg-gray-100 rounded-l-xl rounded-br-xl dark:bg-gray-700">
        <div class="flex items-center justify-end space-x-2 rtl:space-x-reverse">
          <span class="text-sm font-semibold text-gray-900 dark:text-white">Bonnie Green</span>
          <span class="text-sm font-normal text-gray-500 dark:text-gray-400">11:46</span>
        </div>
        <p class="text-sm font-normal py-2.5 text-gray-900 dark:text-white">That's awesome. I think our users will really appreciate the improvements.</p>
        <span class="text-sm font-normal text-gray-500 dark:text-gray-400">Delivered</span>
      </div>
    </div>
    <!-- Chat message 2 (Response) -->
    <div class="flex w-full justify-end p-3">
      <div class="flex flex-col w-full max-w-[320px] leading-1.5 p-4 border-gray-200 bg-gray-100 rounded-l-xl rounded-br-xl dark:bg-gray-700">
        <div class="flex items-center justify-end space-x-2 rtl:space-x-reverse">
          <span class="text-sm font-semibold text-gray-900 dark:text-white">Bonnie Green</span>
          <span class="text-sm font-normal text-gray-500 dark:text-gray-400">11:46</span>
        </div>
        <p class="text-sm font-normal py-2.5 text-gray-900 dark:text-white">That's awesome. I think our users will really appreciate the improvements.</p>
        <span class="text-sm font-normal text-gray-500 dark:text-gray-400">Delivered</span>
      </div>
    </div>
    <!-- Chat message 2 (Response) -->
    <div class="flex w-full justify-end p-3">
      <div class="flex flex-col w-full max-w-[320px] leading-1.5 p-4 border-gray-200 bg-gray-100 rounded-l-xl rounded-br-xl dark:bg-gray-700">
        <div class="flex items-center justify-end space-x-2 rtl:space-x-reverse">
          <span class="text-sm font-semibold text-gray-900 dark:text-white">Bonnie Green</span>
          <span class="text-sm font-normal text-gray-500 dark:text-gray-400">11:46</span>
        </div>
        <p class="text-sm font-normal py-2.5 text-gray-900 dark:text-white">That's awesome. I think our users will really appreciate the improvements.</p>
        <span class="text-sm font-normal text-gray-500 dark:text-gray-400">Delivered</span>
      </div>
    </div>
  </div>
  <div class="flex-none">
    <label for="chat" class="sr-only">Your message</label>
    <div class="flex items-center px-3 py-2 rounded-lg bg-gray-50 dark:bg-gray-700">
      <Textarea bind:value={message} id="chat" class="mx-4 bg-white dark:bg-gray-800" rows="1" placeholder="Your message..." />
      <ToolbarButton onclick={sendMessage} color="blue" class="rounded-full text-primary-600 dark:text-primary-500">
        <PaperPlaneOutline class="w-6 h-6 rotate-45" />
        <span class="sr-only">Send message</span>
      </ToolbarButton>
    </div>
  </div>
</div>

<style lang="less">
</style>
