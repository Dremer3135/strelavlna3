<script lang="ts">
  import { pocketbase } from "$lib/pocketbase";
  import { createEventDispatcher, onMount } from "svelte";
  import { currentUser } from "$lib/stores/auth";
  import type { Contest } from "$lib/types/register";

  const dispatch = createEventDispatcher();

  let { contest }: { contest: Contest } = $props();

  let name = $state("");
  let player1email = $state("");
  let player2email = $state("");
  let player3email = $state("");
  let player4email = $state("");
  let player5email = $state("");
  let player1name = $state("");
  let player2name = $state("");
  let player3name = $state("");
  let player4name = $state("");
  let player5name = $state("");
  let error = $state("");

  onMount(() => {
      if (!contest) {
          error = "Contest object is missing.";
      }
  });

  async function handleAddTeam() {
    error = "";
    if (!name) {
      error = "Jméno týmu je potřeba.";
      return;
    }

    try {
      await pocketbase.collection("teams").create({
        name,
        player1email,
        player2email,
        player3email,
        player4email,
        player5email,
        player1name,
        player2name,
        player3name,
        player4name,
        player5name,
        contest: contest.id,
        teacher: $currentUser?.id,
      });
      dispatch("created");
      dispatch("close");
    } catch (err) {
      console.error("Add Team Error:", err);
      error = "Failed to add team.";
    }
  }
</script>

<div class="modal-backdrop" on:click={() => dispatch("close")}>
  <div class="modal-content" on:click|stopPropagation>
    <h2 class="title">Registrovat tým do:<br> <span class="contest-name">{contest.name}</span></h2>
    <form on:submit|preventDefault={handleAddTeam}>
      <input class="team-name" type="text" placeholder="Název týmu" bind:value={name} required />
      <input class="email1" type="email" placeholder="Email hráče 1" bind:value={player1email} />
      <input class="email1" type="text" placeholder="Jméno hráče 1" bind:value={player1name} />
      <input class="email2" type="email" placeholder="Email hráče 2" bind:value={player2email} />
      <input class="email2" type="text" placeholder="Jméno hráče 2" bind:value={player2name} />
      <input class="email3" type="email" placeholder="Email hráče 3" bind:value={player3email} />
      <input class="email3" type="text" placeholder="Jméno hráče 3" bind:value={player3name} />
      <input class="email4" type="email" placeholder="Email hráče 4" bind:value={player4email} />
      <input class="email4" type="text" placeholder="Jméno hráče 4" bind:value={player4name} />
      <input class="email5" type="email" placeholder="Email hráče 5" bind:value={player5email} />
      <input class="email5" type="text" placeholder="Jméno hráče 5" bind:value={player5name} />
      <button type="submit">Registrovat</button>
    </form>
    {#if error}
      <p class="error">{error}</p>
    {/if}
  </div>
</div>

<style>
  .title {
    font-size: 20px;
    font-weight: 600;
    color: #6a83a0;
    font-family: "Lexend";
  }
  .contest-name {
    color: #1c3047;
    font-weight: 800;
    font-size: 30px;
    font-family: "Fredoka";
  }
  .modal-backdrop {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 102;
  }
  .modal-content {
    background: white;
    /* padding: 2rem; */
    border-radius: 8px;
    color: #333;
    padding: 20px 40px;
  }
  form {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  .error {
    color: red;
  }

  input {
    all: unset;
    padding: 10px 5px;
    width: 300px;
  }

  input:hover {
    background-color: #f8f8f8;
  }

  .team-name:focus {
    outline: 3px dashed #EBAD00;
  }
  .email1:focus {
    outline: 3px dashed #EB6E00;
  }
  .email2:focus {
    outline: 3px dashed #EB0072;
  }
  .email3:focus {
    outline: 3px dashed #9500EB;
  }
  .email4:focus {
    outline: 3px dashed #EBAD00;
  }
  .email5:focus {
    outline: 3px dashed #EB6E00;
  }

</style>
