<script lang="ts">
   import goprivate from "./goprivate.api"
   import ed from "./encdecrypt"
   import emojiHandler from "./emojiHandler"

   let clearTextMsg = ""
   let urlForCreatedNote = ""
   let messageContainsEmojis = false
   let linkCopiedToClipboard = false
   let usingUserDefinedPassword = false
   let userDefinedPassword = ""

   const submitNote = async (e: SubmitEvent) => {
      console.log("got form submit", { event: e })
      if (clearTextMsg === "") return
      if (usingUserDefinedPassword && userDefinedPassword === "") return

      const { encryptedHex, ivHex, keyHex } = usingUserDefinedPassword
         ? await ed.encryptWithPassword(clearTextMsg, userDefinedPassword)
         : await ed.encrypt(clearTextMsg)

      try {
         const res = await goprivate.newNote(`${ivHex}\$${encryptedHex}`)
         console.log({ res })
         clearTextMsg = ""
         const passwordPart = userDefinedPassword ? "" : `#${keyHex}`
         urlForCreatedNote = `${document.location.href}read/${res.noteID}${passwordPart}`
      } catch (error) {
         console.log({ error })
      }
   }

   $: {
      messageContainsEmojis = emojiHandler.containsEmojis(clearTextMsg)
   }
</script>

<article>
   <h1>New Note</h1>
   {#if urlForCreatedNote === ""}
      <form on:submit|preventDefault={submitNote}>
         <textarea
            bind:value={clearTextMsg}
            style="resize: vertical;"
            cols="30"
            rows="10"
         />
         <div>
            <input
               bind:checked={usingUserDefinedPassword}
               class="m-0"
               type="checkbox"
               name="passwordOption"
               id="passwordOption"
            />
            <label class="m-0" for="passwordOption"
               ><b>Set own password</b></label
            >
         </div>
         {#if usingUserDefinedPassword}
            <p>
               <b>NB:</b>When setting your own password it will not be part of
               the generated link. You are expected to share it with the message
               recipient what ever way you see fit.
            </p>
            <label for="userPassword"><b>Password to use:</b></label>
            <input
               bind:value={userDefinedPassword}
               type="text"
               name="userPassword"
               id="userPassword"
            />
         {/if}
         <br />
         {#if messageContainsEmojis}
            <p><b>NB: As of now, your message cannot contain emojis.</b></p>
            <p>Remove them in order to submit the note.</p>
         {/if}
         {#if !messageContainsEmojis}
            <button on:click|stopPropagation type="submit">Create note</button>
         {/if}
      </form>
   {/if}
   <br />
   {#if urlForCreatedNote !== ""}
      <form on:submit|preventDefault={() => {}}>
         <label for="notelink"><h2>Note created!</h2></label>
         <p>Your note can be accessed using this link.</p>
         <input
            on:click={(ev) => {
               const link = ev.currentTarget.value
               console.log({ link })
               navigator.clipboard.writeText(link)
               linkCopiedToClipboard = true
            }}
            id="notelink"
            type="text"
            value={urlForCreatedNote}
         />
         {#if linkCopiedToClipboard}
            <p><mark>Link is copied to clipboard!</mark></p>
         {/if}
      </form>
   {/if}
</article>

<style>
</style>
