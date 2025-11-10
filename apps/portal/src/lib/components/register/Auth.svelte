<script lang="ts">
  import { pocketbase } from "$lib/pocketbase";
  import { createEventDispatcher } from "svelte";

  const dispatch = createEventDispatcher();

  let email = $state("");
  let password = $state("");
  let passwordConfirm = $state("");
  let school = $state("");
  let error = $state("");
  let { type = "login" }: { type?: "login" | "register" } = $props();

  let selectedKraj = $state("");
  let selectedOkres = $state("");
  let schools = $state<any[]>([]);

  const kraje = [
    { name: "Hlavní město Praha", okresy: ["Praha 1", "Praha 2", "Praha 3", "Praha 4", "Praha 5", "Praha 6", "Praha 7", "Praha 8", "Praha 9", "Praha 10", "Praha 11", "Praha 12", "Praha 13", "Praha 14", "Praha 15", "Praha 16", "Praha 17", "Praha 18", "Praha 19", "Praha 20", "Praha 21", "Praha 22"] },
    { name: "Jihočeský kraj", okresy: ["České Budějovice", "Český Krumlov", "Jindřichův Hradec", "Písek", "Prachatice", "Strakonice", "Tábor"] },
    { name: "Jihomoravský kraj", okresy: ["Blansko", "Brno-město", "Brno-venkov", "Břeclav", "Hodonín", "Vyškov", "Znojmo"] },
    { name: "Karlovarský kraj", okresy: ["Cheb", "Karlovy Vary", "Sokolov"] },
    { name: "Královéhradecký kraj", okresy: ["Hradec Králové", "Jičín", "Náchod", "Rychnov nad Kněžnou", "Trutnov"] },
    { name: "Liberecký kraj", okresy: ["Česká Lípa", "Jablonec nad Nisou", "Liberec", "Semily"] },
    { name: "Moravskoslezský kraj", okresy: ["Bruntál", "Frýdek-Místek", "Karviná", "Nový Jičín", "Opava", "Ostrava-město"] },
    { name: "Olomoucký kraj", okresy: ["Jeseník", "Olomouc", "Prostějov", "Přerov", "Šumperk"] },
    { name: "Pardubický kraj", okresy: ["Chrudim", "Pardubice", "Svitavy", "Ústí nad Orlicí"] },
    { name: "Plzeňský kraj", okresy: ["Domažlice", "Klatovy", "Plzeň-město", "Plzeň-jih", "Plzeň-sever", "Rokycany", "Tachov"] },
    { name: "Středočeský kraj", okresy: ["Benešov", "Beroun", "Kladno", "Kolín", "Kutná Hora", "Mělník", "Mladá Boleslav", "Nymburk", "Praha-východ", "Praha-západ", "Příbram", "Rakovník"] },
    { name: "Ústecký kraj", okresy: ["Děčín", "Chomutov", "Litoměřice", "Louny", "Most", "Teplice", "Ústí nad Labem"] },
    { name: "Kraj Vysočina", okresy: ["Havlíčkův Brod", "Jihlava", "Pelhřimov", "Třebíč", "Žďár nad Sázavou"] },
    { name: "Zlínský kraj", okresy: ["Kroměříž", "Uherské Hradiště", "Vsetín", "Zlín"] },
  ];

  async function fetchSchools() {
    if (!selectedKraj || !selectedOkres) {
      schools = [];
      return;
    }

    try {
      const result = await pocketbase.collection("skoly").getFullList({
        filter: `kraj = "${selectedKraj}" && okres = "${selectedOkres}"`,
      });
      schools = result;
    } catch (err) {
      console.error("Failed to fetch schools:", err);
      schools = [];
    }
  }

  async function handleLoginSubmit() {
    error = "";
    try {
      await pocketbase.collection("teachers").authWithPassword(email, password);
      // Auth is successful, authStore will update, App.svelte will react.
      dispatch("close");
    } catch (err: any) {
      if (err.status === 400) {
        error = "Zadané údaje nesedí";
      } else if (err.status === 403) {
        error = "Váš účet zatím není ověřený, ale pracujeme na tom :)";
      }
      else {
        error = "Failed to authenticate. Please check your credentials.";
      }
    }
  }

  async function handleRegisterSubmit() {
    error = "";
    if (password !== passwordConfirm) {
      error = "Passwords do not match.";
      return;
    }

    try {
      console.log("Registering with school ID:", school);
      await pocketbase.collection("teachers").create({
        email,
        password,
        passwordConfirm,
        skola: school,
      });

      // After creating, log them in
      await pocketbase.collection("teachers").authWithPassword(email, password);
      dispatch("close");
    } catch (err: any) {
      console.error("Registration Error:", err);
      error = "Failed to register. Please try again.";
    }
  }
</script>

<div class="auth-modal-backdrop" on:click={() => dispatch("close")}>
  <div class="auth-modal-content" on:click|stopPropagation>
    {#if type === "login"}
      <form on:submit|preventDefault={handleLoginSubmit}>
        <h2>Přihlásit se</h2>
        <input type="email" placeholder="Email" bind:value={email} required class="email"/>
        <input
          type="password"
          placeholder="Heslo"
          bind:value={password}
          required
          class="password"
        />
        <button type="submit">Pokračovat</button>
      </form>
    {:else if type === "register"}
      <form on:submit|preventDefault={handleRegisterSubmit}>
        <h2>Registrovat se</h2>
        <input type="email" bind:value={email} 
          placeholder="email"
          class="email"
          required
        />
        <input
          type="password"
          placeholder="Heslo"
          bind:value={password}
          required
          class="password"
        />
        <input
          type="password"
          placeholder="Potvrďte heslo"
          bind:value={passwordConfirm}
          required
          class="password-again"
        />
        <select bind:value={selectedKraj} class="school" on:change={() => {selectedOkres = ""}}>
          <option value="">Vyberte kraj</option>
          {#each kraje as kraj}
            <option value={kraj.name}>{kraj.name}</option>
          {/each}
        </select>
        <select bind:value={selectedOkres} class="school" on:change={fetchSchools} disabled={!selectedKraj}>
          <option value="">Vyberte okres</option>
          {#if selectedKraj}
            {#each kraje.find(k => k.name === selectedKraj)?.okresy ?? [] as okres}
              <option value={okres}>{okres}</option>
            {/each}
          {/if}
        </select>
        <select bind:value={school} class="school" disabled={!selectedOkres || schools.length === 0}>
          <option value="">Vyberte školu</option>
          {#each schools as s}
            <option value={s.id}>{s.plny_nazev}</option>
          {/each}
        </select>
        <button type="submit">Registrovat</button>
      </form>
    {/if}
    {#if error}
      <p class="error">{error}</p>
    {/if}
  </div>
</div>

<style>
  .auth-modal-backdrop {
    all: unset;
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 100;
  }
  .auth-modal-content {
    background: white;
    padding: 2rem;
    border-radius: 8px;
    color: #333;
    width: 450px;
  }

  h2 {
    margin-top: 0px;
  }

  form {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    padding: 20px 10px;
  }
  .error {
    color: red;
  }

  form input {
    all: unset;
    padding: 10px 5px;
  }
  form input:active {
    border-bottom: 2px black;
  }

  select {
    all: unset;
    padding: 10px 5px;
  }

  .email:hover, .password:hover, .password-again:hover, .school:hover {
    background-color: #f8f8f8;
  }

  .email:focus {
    outline: 3px dashed #EBAD00;
  }
  .password:focus {
    outline: 3px dashed #EB6E00;
  }
  .password-again:focus {
    outline: 3px dashed #EB0072;
  }
  .school:focus {
    outline: 3px dashed #9500EB;
  }
</style>
