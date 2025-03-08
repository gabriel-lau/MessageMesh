<script lang="ts">
  import { Listgroup, ListgroupItem, Badge, Indicator } from 'flowbite-svelte';
  let { selectedPeer = $bindable<string>(), peerList = $bindable<string[]>([]) } = $props();
  function selectPeer(peer: string) {
    selectedPeer = peer;
  }

</script>

<div class="flex-none h-screen">
  <Listgroup active class="h-full w-96 rounded-none">
    <h3 class="p-1 text-center text-xl font-medium text-gray-900 dark:text-white">Chats</h3>
    {#each peerList as peer}
      <ListgroupItem 
        class="text-base font-semibold gap-2 hover:bg-primary-100 dark:hover:bg-gray-900 {selectedPeer === peer ? 'bg-gray-100 dark:bg-gray-900' : ''}" 
        attrs={{
          onclick: () => selectPeer(peer)
        }}
      >
        <div class="flex items-center place-content-between w-full">
          <div class="flex-1 min-w-0">
            <p class="text-sm font-semibold text-gray-900 truncate dark:text-white">{peer}</p>
          </div>
          <Badge color="green" rounded class="px-2.5 py-0.5">
            <Indicator color="green" size="xs" class="me-1" />Available
          </Badge>
        </div>
      </ListgroupItem>
    {/each}
  </Listgroup>
</div>

<style>
</style>
