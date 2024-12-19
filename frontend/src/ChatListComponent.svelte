<script lang="ts">
  export let isChatOpen;
</script>

<div class={isChatOpen ? 'chatlistcomponent-wrapper d-none d-lg-block' : ''}>ChatList</div>

<style>
  .chatlistcomponent-wrapper {
    height: 100vh;
    background-color: var(--md-sys-color-surface-container-low);
    width: 400px;
  }
</style>
