/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    // path to templates dir inside docker container. When setting up locally, we will need to cange the path.
    "/templates/**/*.{gohtml,html}" // include files in templates dir and its subdir having extension gohtml & html
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}

