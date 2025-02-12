<script lang="ts">
  import './app.css';
  import NavigationRailComponent from './components/NavigationRailComponent.svelte';
  import ChatListComponent from './components/ChatListComponent.svelte';
  import ChatComponent from './components/ChatComponent.svelte';
  import * as Wails from '../wailsjs/runtime/runtime.js';
  import { models } from '../wailsjs/go/models.js';

  let userPeerID = $state('');
  Wails.EventsOn("getUserPeerID", (data: string) => {
    userPeerID = data;
  });
  let blockchain = $state<models.Block[]>([]);
  Wails.EventsOn("getBlockchain", (data: models.Block[]) => {
    blockchain = data;
  });
  let messages = $state<models.Message[]>([]);
  Wails.EventsOn("getMessages", (data: models.Message[]) => {
    messages = data;
  });
  let accounts = $state<models.Account[]>([]);
  Wails.EventsOn("getAccounts", (data: models.Account[]) => {
    accounts = data;
  });

  let selectedPeer = $state('');
</script>

<main>
  <div class="flex w-screen h-screen bg-primary-50 dark:bg-gray-900">
    <NavigationRailComponent></NavigationRailComponent>
    <div class="flex flex-row w-full">
      <ChatListComponent bind:userPeerID bind:selectedPeer bind:accounts></ChatListComponent>
      <ChatComponent bind:userPeerID bind:selectedPeer bind:messages></ChatComponent>
    </div>
  </div>
  <!-- <img alt="Wails logo" id="logo" src="{logo}"> -->
  <!-- <div class="result" id="result">{resultText}</div>
  <div class="input-box" id="input">
    <input autocomplete="off" bind:value={name} class="input" id="name" type="text"/>
    <button class="btn" on:click={greet}>Greet</button>
  </div> -->
</main>

<style>
</style>
