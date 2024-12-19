
<script>
    import '@material/web/fab/fab';
</script>

<div class="naviagation-wrapper d-flex flex-column">
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
</div>

<style>
    .naviagation-wrapper {
        width: 80px;
        height: 100vh;
        background-color: var(--md-sys-color-surface-container-lowest);
    }

    .fab-wrapper {
        width: 100%;
        display: flex;
        justify-content: center;
        align-items: center;
    }

    .navigation-options {
        width: 100%;
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
    }

    .navigation-item {
        width: 100%;
        display: flex;
        justify-content: center;
        align-items: center;
    }
</style>