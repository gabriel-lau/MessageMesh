<script lang="ts">
  import {
    Sidebar,
    SidebarGroup,
    SidebarItem,
    SidebarWrapper,
    Indicator,
  } from 'flowbite-svelte';
  import {
    DrawSquareOutline,
  } from 'flowbite-svelte-icons';
  let { onlinePeerList = $bindable<string[]>([]), online = $bindable<boolean>(false) } = $props();

  let connected: 'Connected' | 'Disconnected' = $state('Disconnected');
  let status: 'green' | 'red' = $state('red');
  $effect(() => {
    if (online) {
      status = 'green';
      connected = 'Connected';
    } else {
      status = 'red';
      connected = 'Disconnected';
    }
  });

  let spanClass = 'text-xs text-center';
</script>

<div class="h-screen">
  <Sidebar class="w-auto h-full">
    <SidebarWrapper class="h-full rounded-none">
      <SidebarGroup>
        <SidebarItem label="Connected to {onlinePeerList.length} peers" class="flex flex-col" {spanClass} href="#">
          <svelte:fragment slot="icon">
            <DrawSquareOutline class="w-6 h-6 text-gray-500 transition duration-75 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white" />
          </svelte:fragment>
        </SidebarItem>

        <SidebarItem label={connected} class="flex flex-col" {spanClass} href="#">
          <svelte:fragment slot="icon">
            <Indicator class="m-2" color={status} />
          </svelte:fragment>
        </SidebarItem>
      </SidebarGroup>
    </SidebarWrapper>
  </Sidebar>
</div>

<style>
</style>
