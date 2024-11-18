<script>
	export let isChatOpen;
    import '@material/web/textfield/outlined-text-field';
    import '@material/web/iconbutton/filled-icon-button';
</script>

<div class={isChatOpen ? " d-block chatcomponent": "chatcomponent"}>
    <div class="chatcomponent-wrapper d-flex flex-column">
        <div class="topbar">

        </div>
        <div class="chat flex-grow-1">
    
        </div>
        <div class="bottombar">
            <div class="bottombar-wrapper d-flex flex-row align-content-center align-items-center">
                <md-outlined-text-field class="send-field m-3" placeholder="Send">
                </md-outlined-text-field>
                <md-filled-icon-button aria-label="Copy" class="m-3">
                    <md-icon class="material-icons">send</md-icon>
                </md-filled-icon-button>
            </div>
        </div>    
    </div>
</div>

<style lang=less>
    .chatcomponent {
        width: 100%;
        height: 100vh;
        background-color: var(--md-sys-color-surface-container-highest);
    }

    .chatcomponent-wrapper {
        height: 100vh;
        background-color: var(--md-sys-color-surface-container-highest);
    }

    .topbar {
        height: 64px;
        background-color: var(--md-sys-color-surface-container);
    }

    .chat {
        
    }

    .bottombar {
        background-color: var(--md-sys-color-surface-container);
        .send-field {
            width: 100%;
        }
    }
</style>