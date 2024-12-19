<script lang="ts">
  import { Listgroup, ListgroupItem, Badge, Indicator } from 'flowbite-svelte';
  export let isChatOpen;
</script>

<div class="flex-none h-screen">
  <Listgroup active class="h-full w-96">
    <h3 class="p-1 text-center text-xl font-medium text-gray-900 dark:text-white">User list</h3>
    <ListgroupItem class="text-base font-semibold gap-2">
      <div class="flex items-center place-content-between w-full">
        <div class="flex-1 min-w-0">
          <p class="text-sm font-semibold text-gray-900 truncate dark:text-white">Neil Sims</p>
          <p class="text-sm text-gray-500 truncate dark:text-gray-400">email@flowbite.com</p>
        </div>
        <Badge color="green" rounded class="px-2.5 py-0.5">
          <Indicator color="green" size="xs" class="me-1" />Available
        </Badge>
      </div>
    </ListgroupItem>
    <ListgroupItem class="text-base font-semibold gap-2">
      <div class="flex items-center place-content-between w-full">
        <div class="flex-1 min-w-0">
          <p class="text-sm font-semibold text-gray-900 truncate dark:text-white">Neil Sims</p>
          <p class="text-sm text-gray-500 truncate dark:text-gray-400">email@flowbite.com</p>
        </div>
        <Badge color="green" rounded class="px-2.5 py-0.5">
          <Indicator color="green" size="xs" class="me-1" />Available
        </Badge>
      </div>
    </ListgroupItem>
    <ListgroupItem class="text-base font-semibold gap-2">
      <div class="flex items-center place-content-between w-full">
        <div class="flex-1 min-w-0">
          <p class="text-sm font-semibold text-gray-900 truncate dark:text-white">Neil Sims</p>
          <p class="text-sm text-gray-500 truncate dark:text-gray-400">email@flowbite.com</p>
        </div>
        <Badge color="green" rounded class="px-2.5 py-0.5">
          <Indicator color="green" size="xs" class="me-1" />Available
        </Badge>
      </div>
    </ListgroupItem>
  </Listgroup>
</div>

<style>
</style>
