<script lang="ts">
  import logo from './assets/images/logo-universal.png'
  import {Greet} from '../wailsjs/go/main/App.js'
  import NavigationRail from './NavigationRailComponent.svelte';
  import ChatListComponent from './ChatListComponent.svelte';
  import ChatComponent from './ChatComponent.svelte';

  let name2 = $state('');

  let resultText: string = "Please enter your name below 👇"
  let name: string
  let isChatOpen= true;


  function greet(): void {
    Greet(name).then(result => resultText = result)
  }
</script>

<main>
  <div class="main-wrapper dark">
    <NavigationRail/>
    <div class="d-inline-flex" style="width: 100%;">
      <ChatListComponent bind:isChatOpen={isChatOpen}/>
      <ChatComponent bind:isChatOpen={isChatOpen}/>  
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
  @import './css/dark.css';

  .main-wrapper {
    display: flex;
    height: 100vh;
    /* background-color: var(--md-sys-color-surface-container-lowest); */
  }
</style>
