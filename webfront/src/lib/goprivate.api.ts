let base = "http://localhost:8080"

if (import.meta.env.MODE === "production") {
   base = ""
}

const apiBase = `${base}/api/v1`

console.log({ apiBase })

export default {
   /* 
    Ping
  */
   ping: async function () {
      const resp = await fetch(`${apiBase}/ping`, {
         method: "GET",
         headers: {
            // Authorization: `Bearer ${getToken()}`
         }
      })
      if (resp.status !== 200) {
         throw resp
      }
      const jsonRes = await resp.json()
      return jsonRes
   },

   /* 
   New note
   */
   newNote: async function (
      clearTextContent: string
   ): Promise<{ noteID: string; key: string }> {
      /* 
    [ ] NEEDS to encrypt note content
    */

      const resp = await fetch(`${apiBase}/newnote`, {
         method: "POST",
         headers: {
            "Content-Type": "application/json"
         },
         body: JSON.stringify({
            noteContent: clearTextContent
         })
      })
      if (resp.status !== 200) {
         throw resp
      }
      let json = await resp.json()
      json["key"] = "KEY_FOR_ENCRYPTION"
      return json
   }
}
