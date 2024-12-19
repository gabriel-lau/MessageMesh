<script lang="ts">
  import { Sidebar, SidebarBrand, SidebarCta, SidebarDropdownItem, SidebarDropdownWrapper, SidebarGroup, SidebarItem, SidebarWrapper } from 'flowbite-svelte';
  import { ChartPieSolid, GridSolid, MailBoxSolid, UserSolid, ArrowRightToBracketOutline, EditOutline } from 'flowbite-svelte-icons';
  let spanClass = 'flex-1 ms-3 whitespace-nowrap';
</script>

<div class="h-screen">
  <Sidebar class="h-full w-auto">
    <SidebarWrapper>
      <SidebarGroup>
        <SidebarItem>
          <svelte:fragment slot="icon">
            <ChartPieSolid class="w-6 h-6 text-gray-500 transition duration-75 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white" />
          </svelte:fragment>
        </SidebarItem>
        <SidebarItem {spanClass}>
          <svelte:fragment slot="icon">
            <GridSolid class="w-6 h-6 text-gray-500 transition duration-75 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white" />
          </svelte:fragment>
        </SidebarItem>
        <SidebarItem {spanClass}>
          <svelte:fragment slot="icon">
            <MailBoxSolid class="w-6 h-6 text-gray-500 transition duration-75 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white" />
          </svelte:fragment>
        </SidebarItem>
        <SidebarItem href="/docs/components/sidebar">
          <svelte:fragment slot="icon">
            <UserSolid class="w-6 h-6 text-gray-500 transition duration-75 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white" />
          </svelte:fragment>
        </SidebarItem>
        <SidebarItem>
          <svelte:fragment slot="icon">
            <ArrowRightToBracketOutline
              class="w-6 h-6 text-gray-500 transition duration-75 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white"
            />
          </svelte:fragment>
        </SidebarItem>
        <SidebarItem>
          <svelte:fragment slot="icon">
            <EditOutline class="w-6 h-6 text-gray-500 transition duration-75 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white" />
          </svelte:fragment>
        </SidebarItem>
      </SidebarGroup>
    </SidebarWrapper>
  </Sidebar>
</div>

<!-- <div class="naviagation-wrapper d-flex flex-column">
    <div class="fab-wrapper mb-4 mt-4">
        <md-fab aria-label="Edit">
            <md-icon class="material-icons" slot="icon">add</md-icon>
        </md-fab>
    </div>
    <div class="navigation-options">
        <div class="navigation-item">
            ITEM 1
        </div>
        <div class="navigation-item">
            ITEM 2
        </div>
        <div class="navigation-item">
            ITEM 3
        </div>

    </div>
</div> -->

<style>
</style>
