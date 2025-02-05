<script lang="ts">
  import './app.css';
  import logo from './assets/images/logo-universal.png';
  import { Greet } from '../wailsjs/go/main/App.js';
  import NavigationRailComponent from './components/NavigationRailComponent.svelte';
  import ChatListComponent from './components/ChatListComponent.svelte';
  import ChatComponent from './components/ChatComponent.svelte';
  import * as Wails from '../wailsjs/runtime/runtime.js';
  import * as Main from '../wailsjs/go/main/app.js';

  let userPeerID = $state('');
  Wails.EventsOn("getUserPeerID", (peerID: string) => {
    userPeerID = peerID;
  });
</script>

<main>
  <div class="flex w-screen h-screen bg-primary-50 dark:bg-gray-900">
    <NavigationRailComponent></NavigationRailComponent>
    <div class="flex flex-row w-full">
      <ChatListComponent bind:userPeerID></ChatListComponent>
      <ChatComponent bind:userPeerID></ChatComponent>
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
