<script lang="ts">
  import { Sidebar, SidebarGroup, SidebarItem, SidebarWrapper, Indicator, Button, Modal, DarkMode, Listgroup, ListgroupItem } from 'flowbite-svelte';
  import {
    PlusOutline,
    DrawSquareOutline,
    AdjustmentsVerticalOutline,
    CheckCircleSolid,
    ArrowRightToBracketOutline,
    DownloadOutline
  } from 'flowbite-svelte-icons';
  let spanClass = 'text-xs text-center';
  let newChatModal = false;
  let settingsModal = false;
  let settingsToast = false;
</script>

<div class="h-screen">
  <Sidebar class="w-auto h-full">
    <SidebarWrapper class="h-full">
      <SidebarGroup>
        <SidebarItem class="flex flex-col" {spanClass} on:click={() => (newChatModal = true)} href="#">
          <svelte:fragment slot="icon">
            <Button class="!p-2"><PlusOutline class="w-6 h-6" /></Button>
          </svelte:fragment>
        </SidebarItem>

        <SidebarItem label="20 Peers Online" class="flex flex-col" {spanClass} href="#">
          <svelte:fragment slot="icon">
            <DrawSquareOutline class="w-6 h-6 text-gray-500 transition duration-75 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white" />
          </svelte:fragment>
        </SidebarItem>

        <SidebarItem label="Online" class="flex flex-col" {spanClass} href="#">
          <svelte:fragment slot="icon">
            <Indicator class="m-2" color="green" />
          </svelte:fragment>
        </SidebarItem>

        <SidebarItem label="Settings" class="flex flex-col" {spanClass} on:click={() => (settingsModal = true)} href="#">
          <svelte:fragment slot="icon">
            <AdjustmentsVerticalOutline
              class="w-6 h-6 text-gray-500 transition duration-75 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white"
            />
          </svelte:fragment>
        </SidebarItem>
      </SidebarGroup>
    </SidebarWrapper>
  </Sidebar>
  <Modal title="Settings" size="xs" bind:open={settingsModal} autoclose>
    <Listgroup active class="w-full" on:click={console.log}>
      <ListgroupItem>
        <a href="#" class="flex items-center gap-2">
          <ArrowRightToBracketOutline class="w-5 h-5" />
          <span>Logout</span>
        </a>
      </ListgroupItem>
      <ListgroupItem>
        <a href="#" class="flex items-center gap-2">
          <DownloadOutline class="w-5 h-5" />
          <span>Download</span>
        </a>
      </ListgroupItem>
    </Listgroup>
    <DarkMode />

    <svelte:fragment slot="footer">
      <Button color="alternative">Close</Button>
    </svelte:fragment>
  </Modal>
  <Modal title="New Chat" bind:open={newChatModal} autoclose>
    <p class="text-base leading-relaxed text-gray-500 dark:text-gray-400">
      The European Union’s General Data Protection Regulation (G.D.P.R.) goes into effect on May 25 and is meant to ensure a common set of data rights in the
      European Union. It requires organizations to notify users as soon as possible of high-risk data breaches that could personally affect them.
    </p>
    <svelte:fragment slot="footer">
      <Button on:click={() => alert('Handle "success"')}>Save</Button>
      <Button color="alternative">Cancel</Button>
    </svelte:fragment>
  </Modal>
  <!-- <Toast dismissable={settingsToast} color="green">
    <svelte:fragment slot="icon">
      <CheckCircleSolid class="w-5 h-5" />
      <span class="sr-only">Check icon</span>
    </svelte:fragment>
    Item moved successfully.
  </Toast> -->
</div>

<style>
</style>
