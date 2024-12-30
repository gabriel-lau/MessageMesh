<script lang="ts">
  import {
    Sidebar,
    SidebarGroup,
    SidebarItem,
    SidebarWrapper,
    Indicator,
    Button,
    Modal,
    DarkMode,
    Listgroup,
    ListgroupItem,
    Input,
    Helper,
    ButtonGroup,
    InputAddon
  } from 'flowbite-svelte';
  import {
    PlusOutline,
    DrawSquareOutline,
    AdjustmentsVerticalOutline,
    ArrowRightToBracketOutline,
    DownloadOutline,
    UserCircleSolid
  } from 'flowbite-svelte-icons';
  import { EventsOn } from '../../wailsjs/runtime/runtime';
  import * as Wails from '../../wailsjs/runtime/runtime.js';
  let peerList = $state([]);
  Wails.EventsOn('getPeerList', (newPeerList) => {
    console.log(newPeerList);
    peerList = newPeerList;
  });

  let spanClass = 'text-xs text-center';
  let newChatModal = false;
  let status: 'green' | 'red' | 'disabled' | 'gray' = 'green';
  let settingsModal = false;
  let settingsToast = false;
  // function getOnlinePeers() {
  //   const eventRemover = EventsOn('myEvent', () => console.log('Event fired.'));
  //   return () => eventRemover();
  // }
</script>

<div class="h-screen">
  <Sidebar class="w-auto h-full">
    <SidebarWrapper class="h-full rounded-none">
      <SidebarGroup>
        <SidebarItem class="flex flex-col" {spanClass} on:click={() => (newChatModal = true)} href="#">
          <svelte:fragment slot="icon">
            <Button class="!p-2"><PlusOutline class="w-6 h-6" /></Button>
          </svelte:fragment>
        </SidebarItem>

        <SidebarItem label="Connected to {peerList.length} peers" class="flex flex-col" {spanClass} href="#">
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
          <span>Download Keys</span>
        </a>
      </ListgroupItem>
    </Listgroup>

    <svelte:fragment slot="footer">
      <div class="flex justify-between">
        <Button color="alternative">Close</Button>
        <DarkMode />
      </div>
    </svelte:fragment>
  </Modal>
  <Modal title="New Chat" size="sm" bind:open={newChatModal} autoclose>
    <div class="mb-6">
      <ButtonGroup class="w-full">
        <InputAddon>
          <UserCircleSolid class="w-4 h-4 text-gray-500 dark:text-gray-400" />
        </InputAddon>
        <Input color={status} id="website-admin" placeholder="Enter account username" />
      </ButtonGroup>
      <Helper class="mt-2" color={status}>
        <span class="font-medium">Well done!</span>
        Some success message.
      </Helper>
    </div>
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
