<script lang="ts">
    import logo from "$lib/assets/images/logo.svg";
    import for_teachers from "$lib/assets/images/navbar/for_teachers.svg"
    import { pocketbase } from "$lib/pocketbase";
    import { currentUser } from "$lib/stores/auth";
    import { onMount, createEventDispatcher } from "svelte";
    import { myTeamsStore } from "$lib/stores/myTeams";
    import { page } from "$app/state";

    let register_hovered = false;
    let login_hovered = false;

    const dispatch = createEventDispatcher();

    function showAuthLogin() {
        dispatch('showAuth', "login");
    }
    function showAuthRegister() {
        dispatch('showAuth', "register");
    }

    onMount(() => {
        console.log("cyurr: ", $currentUser);
    })

</script>

<main class:logged-in={$currentUser}>
    <a href="/">
        <img src={logo} alt="Strela Vlna Logo" class="logo">
    </a>
    
    <a href="/home/about_us" class="nav-element nav2">O nás</a>
    <a href="/home/current-season" class="nav-element nav3">Aktuální ročník</a>
    <a href="/home/archive" class="nav-element nav4">Archiv</a>
    <a href="/home/rules" class="nav-element nav5">Pravidla</a>
    <div class="right-section">
        {#if $currentUser}
        <a href="/logout" class="logout" data-sveltekit-reload>
            <p>Odhlásit se</p>
            <i class="fa-solid fa-arrow-right-from-bracket"></i>
        </a>
        {:else}
            <img src={for_teachers} alt="For Teachers">
            <div class="button-holder">
                <a class="register" href="/register?redirectTo={page.url.pathname}">Registrovat</a>
                <div class="bottom_row">
                    <p class="or">nebo</p>
                    <a class="login" href="/login?redirectTo={page.url.pathname}">Přihlásit</a>
                </div>
            </div>
        {/if}
    </div>
</main>

<style lang="scss">
    main {
        height: 120px;
        background-image: linear-gradient(to right, var(--color-lightblue) 0%, var(--color-blue) 100%);
        width: 100%;
        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: space-between;
        padding-inline: 50px;
        box-sizing: border-box;
        z-index: 10;
        border-bottom: 2px dotted var(--color-lightblue);
    }

    a.register, a.login {
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
    a.register:hover, a.login:hover {
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

    .nav-element {
        all: unset;
        cursor: pointer;
        text-wrap: nowrap;
        padding: 10px 15px;
        color: white;
        font-family: 'Fredoka';
        font-size: 20px;
        font-weight: 600;
        transition: all cubic-bezier(0.215, 0.610, 0.355, 1) 0.3s;
    
        &:hover {
            transform: translateY(-5px);
        }
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

    .logo {
        height: 55px;
    }

    .right-section {
        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: center;
        gap: 10px;
    }

    a.logout {
        position: relative;
        all: unset;
        cursor: pointer;
        box-sizing: border-box;
        padding: 10px 20px;
        width: 150px;
        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: center;  
        gap: 0px;      
        
        p {
            font-family: 'Lexend';
            font-size: 16px;
            color: white;
            font-weight: 300;
            margin: 0px;
            transition: all cubic-bezier(0.215, 0.610, 0.355, 1) 0.3s;
        }

        i {
            transform: translateX(0px);
            transition: all cubic-bezier(0.215, 0.610, 0.355, 1) 0.3s;
            color: white;
            opacity: 0;
            // font-weight: 200;
        }

        &:hover {
            p {
                transform: translateX(-5px);
            }

            i {
                transform: translateX(5px);
                opacity: 1;
            }
        }
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