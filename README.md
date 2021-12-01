Go Bar is just a bar I programmed for my Linux Desktop, because I kept getting a lot of bars that didn't have what I wanted.

<img src="./bar.png" alt="Go Bar" width="2560" height="26" />

--------
Add the Fonts to the folder at: $HOME/.config/i3/Github
Create the symbolic links for all the fonts in the Github folder as well inside the ~/.fonts folder:


mkdir Github
cd Github

# Material Design Icons
git clone --depth 1 https://github.com/google/material-design-icons
ln -s $PWD/material-design-icons/font/MaterialIcons-Regular.ttf ~/.fonts/

# Community Fork
git clone --depth 1 https://github.com/Templarian/MaterialDesign-Webfont
ln -s $PWD/MaterialDesign-Webfont/fonts/materialdesignicons-webfont.ttf ~/.fonts/

# FontAwesome
git clone --depth 1 https://github.com/FortAwesome/Font-Awesome
ln -s "$PWD/Font-Awesome/otfs/Font Awesome 5 Free-Solid-900.otf" ~/.fonts/
ln -s "$PWD/Font-Awesome/otfs/Font Awesome 5 Free-Regular-400.otf" ~/.fonts/
ln -s "$PWD/Font-Awesome/otfs/Font Awesome 5 Brands-Regular-400.otf" ~/.fonts/

# Typicons
git clone --depth 1 https://github.com/stephenhutchings/typicons.font
ln -s $PWD/typicons.font/src/font/typicons.ttf ~/.fonts/

# May need to rebuild the font cache using "sudo fc-cache -f -v", and restart i3 to pick up the new fonts.



###### Editing The config.yaml file #######
Create a config.yaml File inside of $HOME/.config/i3/gobar/
add the following:
---
  openweather:
  	key: ""
	cityid: ""

### Adding the GoBar to run with i3 ###
In order for it to run in i3 you just have to make sure the gobar folder is in $HOME/.config/i3/
and add:
status_command exec $HOME/.config/i3/gobar/gobar
