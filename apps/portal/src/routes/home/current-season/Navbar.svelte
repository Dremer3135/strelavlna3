<script lang="ts">
    import logo from "../assets/strelavlna.svg";
    import for_teachers from "../assets/for_teachers.svg"
    import { pb } from "../lib/pocketbase";
    import { onMount, createEventDispatcher } from "svelte";
    import TeamDetailsModal from "./TeamDetailsModal.svelte";
    import { teamsStore } from "../lib/stores";

    let register_hovered = false;
    let login_hovered = false;

    export let currentUser;
    const dispatch = createEventDispatcher();

    function showAuthLogin() {
        dispatch('showAuth', "login");
    }
    function showAuthRegister() {
        dispatch('showAuth', "register");
    }

    function logout() {
        pb.authStore.clear();
        teamsStore.set([]);
    }

    function toggleTeamsPanel() {
        dispatch('toggleTeamsPanel');
    }

    onMount(() => {
        console.log(currentUser);
    })

</script>

<main class:logged-in={currentUser && pb.authStore.isValid && currentUser?.verified}>
    <a href="/">
        <img src={logo} alt="Strela Vlna Logo" class="logo">
    </a>
    <div class="right-section">
        {#if currentUser && pb.authStore.isValid && currentUser?.verified}
            <button class="logout" on:click={logout}>Logout</button>
        {:else}
            <img src={for_teachers} alt="For Teachers">
            <div class="button-holder">
                <button class="register" on:click={showAuthRegister}>Registrovat</button>
                <div class="bottom_row">
                    <p class="or">nebo</p>
                    <button class="login" on:click={showAuthLogin}>Přihlásit</button>
                </div>
            </div>
        {/if}
    </div>
</main>

<style>
    main {
        height: 120px;
        background-image: linear-gradient(to right, #003F88 0%, #33008B 100%);
        width: 100%;
        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: space-between;
        padding-inline: 50px;
        box-sizing: border-box;
    }

    main.logged-in {
        background-image: linear-gradient(to right, #33008B 0%, #3c008b 100%);
    }

    button.register, button.login {
        all: unset;
        text-align: center;
        max-width: 115px;
        font-family: 'Lexend';
        font-weight: 500;
        font-size: 16px;
        color: white;
        cursor: pointer;
        padding: 2px 4px;
        border-radius: 4px;
        transition: background-color 0.2s ease-out, color 0.2s ease-out;
    }
    button.register:hover, button.login:hover {
        background-color: white;
        color: #1e0052;
        transition: background-color 0s, color 0s;
    }
        
    .or {
        all: unset;
        text-align: center;
        font-family: 'Lexend';
        font-weight: 500;
        font-size: 16px;
        color: #9FA6D4;
    }

    button.my_teams {
        all: unset;
        text-align: center;
        max-width: 115px;
        font-family: 'Lexend';
        font-weight: bold;
        font-size: 16px;
        color: white;
        cursor: pointer;
    }

    .right-section {
        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: center;
        gap: 10px;
    }

    button.logout {
        all: unset;
        font-family: 'Lexend';
        font-weight: bold;
        font-size: 16px;
        color: white;
        cursor: pointer;
        margin-left: 2rem;
    }

    .button-holder {
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
        gap: 0px;
    }
    .bottom_row {
        display: flex;
        align-items: center;
        justify-content: center;
        flex-direction: row;
        gap: 5px;
    }
    @media (max-width: 1100px){
        main{
            flex-direction: column;
            align-items: center;
            justify-content: center;
            gap: 10px;
            padding-top: 30px;
            padding-bottom: 30px;
            height: auto;
        }
    }
    @media (max-width: 600px){
        .logo, .right-section{
            max-width: 80vw;
        }
    }

</style>