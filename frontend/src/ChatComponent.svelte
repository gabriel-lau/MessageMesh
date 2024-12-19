<script lang="ts">
  import { Textarea, ToolbarButton } from 'flowbite-svelte';
  import { Navbar, NavBrand, NavLi, NavUl, NavHamburger } from 'flowbite-svelte';
  import { PaperPlaneOutline } from 'flowbite-svelte-icons';

  let { isChatOpen = $bindable() } = $props();
  let name = $state('');

  function sendMessage() {
    console.log('send message');
  }
</script>

<div class={isChatOpen ? ' d-block chatcomponent' : 'chatcomponent'}>
  <div class="chatcomponent-wrapper d-flex flex-column">
    <div class="topbar"></div>
    <div class="chat flex-grow-1"></div>
    <div class="bottombar">
      <form>
        <label for="chat" class="sr-only">Your message</label>
        <div class="flex items-center px-3 py-2 rounded-lg bg-gray-50 dark:bg-gray-700">
          <Textarea id="chat" class="mx-4 bg-white dark:bg-gray-800" rows="1" placeholder="Your message..." />
          <ToolbarButton onclick={sendMessage} type="submit" color="blue" class="rounded-full text-primary-600 dark:text-primary-500">
            <PaperPlaneOutline class="w-6 h-6 rotate-45" />
            <span class="sr-only">Send message</span>
          </ToolbarButton>
        </div>
      </form>
    </div>
  </div>
</div>

<style lang="less">
  .chatcomponent {
    width: 100%;
    height: 100vh;
    background-color: var(--md-sys-color-surface-container-highest);
  }

  .chatcomponent-wrapper {
    height: 100vh;
    background-color: var(--md-sys-color-surface-container-highest);
  }

  .topbar {
    height: 64px;
    background-color: var(--md-sys-color-surface-container);
  }

  .chat {
  }

  .bottombar {
    background-color: var(--md-sys-color-surface-container);
    .send-field {
      width: 100%;
    }
  }
</style>
