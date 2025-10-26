<script lang="ts">
    import type { CorrectorsResponse } from "$lib/pocketbase-types";
    import logo from "$lib/assets/images/logo_small.svg";
    import logout_icon from "$lib/assets/images/icons/logout.svg"
    import Button from "../general/Button.svelte";

    let { user, isAdmin }: { user: CorrectorsResponse | undefined, isAdmin: boolean } = $props();
</script>

<header>
    <div class="background" class:admin={isAdmin}>
    </div>
    <div class="content">
        <div class="title-wrapper">
            <img src={logo} height="100" class="logo">
        </div>
        <div>
            {#if user}
                <span class="welcome-message">Welcome, {user.username}</span>
            {/if}
        </div>
        {#if user}
            <form method="POST" action="/logout">
                <button class="logout" type="submit">Logout</button>
            </form>
        {/if}
    </div>
</header>

<style>
    .background {
        position: absolute;
        width: 100%;
        height: 100%;
        background-color: black;
        z-index: -1;
    }
    .background.admin {
        /* to be done */
    }

    .welcome-message {
        font-size: 15px;
        color: #6e6e6e;
        font-weight: 300;
    }



    .content {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 1rem;
        padding-inline: 50px;
        border-bottom: 1px solid #ccc;
        font-family: 'Lexend';
        box-sizing: border-box;
        height: 100%;
        
    }
    header {
        position:relative;
        box-sizing: border-box;
        height: 120px;
    }

    .title-wrapper {
        display: flex;
        align-items: flex-end;
        justify-content: center;
        gap: 35px;
    }

    img.logo {
        transform: translateY(-5px);
    }


    button.logout {
        all: unset;
        max-height: 30px;
        font-family: 'Lexend';
        font-size: 15px;
        font-weight: 400;
        color: black;
        cursor: pointer;
        padding: 5px 20px;
        background-color: white;
        border-radius: 2px;
        text-align: center;
    }

    button.logout:hover {
        outline: 1px #AAAAAA solid;
        outline-offset: 5px;
    }

</style>