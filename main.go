package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/kardianos/osext"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/zupzup/calories/command"
	"github.com/zupzup/calories/datasource"
	"github.com/zupzup/calories/renderer"
)

// VERSION indicates the version of the binary
const VERSION = "1.0.0"

const configFile string = ".caloriesconf"
const defaultDBFile string = "calories.db"

// A flagset for all subcommands "e.g.: calories config
var commandFlag = flag.NewFlagSet("", flag.ExitOnError)

var (
	weightFlag        float64
	heightFlag        float64
	activityFlag      float64
	birthDayFlag      string
	genderFlag        string
	unitFlag          string
	dateFlag          string
	yesFlag           bool
	commandOutputFlag string
	positionFlag      int
	fileFlag          string

	defaultDateFlag string
	weekFlag        bool
	monthFlag       bool
	histFlag        int
	commandsFlag    bool
	outputFlag      string
	versionFlag     bool
)

func init() {
	commandFlag.Float64Var(&weightFlag, "weight", -1, "your weight")
	commandFlag.Float64Var(&weightFlag, "w", -1, "your weight (shorthand)")
	commandFlag.Float64Var(&heightFlag, "height", -1, "your height")
	commandFlag.Float64Var(&heightFlag, "h", -1, "your height (shorthand)")
	commandFlag.Float64Var(&activityFlag, "activity", -1, "your activity multiplier")
	commandFlag.Float64Var(&activityFlag, "a", -1, "your activity multiplier (shorthand)")
	commandFlag.StringVar(&birthDayFlag, "birthday", "", "your birthday (dd.mm.yyyy)")
	commandFlag.StringVar(&birthDayFlag, "b", "", "your birthday (dd.mm.yyyy) (shorthand)")
	commandFlag.StringVar(&genderFlag, "gender", "male", "your gender")
	commandFlag.StringVar(&genderFlag, "g", "male", "your gender (shorthand)")
	commandFlag.StringVar(&unitFlag, "unit", "metric", "your preferred unit system (metric | imperial)")
	commandFlag.StringVar(&unitFlag, "u", "metric", "your preferred unit system (metric | imperial) (shorthand)")
	commandFlag.StringVar(&dateFlag, "date", "", "date to add an entry on")
	commandFlag.StringVar(&dateFlag, "d", "", "date to add an entry on (shorthand)")
	commandFlag.BoolVar(&yesFlag, "yes", false, "skip confirmations")
	commandFlag.BoolVar(&yesFlag, "y", false, "skip confirmations (shorthand)")
	commandFlag.StringVar(&commandOutputFlag, "output", "terminal", "output format (terminal | json)")
	commandFlag.StringVar(&commandOutputFlag, "o", "terminal", "output format (terminal | json) (shorthand)")
	commandFlag.IntVar(&positionFlag, "position", -1, "position of the entry to clear (1-n)")
	commandFlag.IntVar(&positionFlag, "p", -1, "position of the entry to clear (1-n) (shorthand)")
	commandFlag.StringVar(&fileFlag, "file", "", "file to export to / import from")
	commandFlag.StringVar(&fileFlag, "f", "", "file to export to / import from (shorthand)")

	flag.StringVar(&defaultDateFlag, "date", "", "date to show")
	flag.StringVar(&defaultDateFlag, "d", "", "date to show (shorthand)")
	flag.BoolVar(&weekFlag, "week", false, "show current week")
	flag.BoolVar(&weekFlag, "w", false, "show current week (shorthand)")
	flag.BoolVar(&monthFlag, "month", false, "show current month")
	flag.BoolVar(&monthFlag, "m", false, "show current month (shorthand)")
	flag.IntVar(&histFlag, "hist", 0, "length of history to show")
	flag.IntVar(&histFlag, "h", 0, "length of history to show (shorthand)")
	flag.BoolVar(&commandsFlag, "commands", false, "show list of commands")
	flag.BoolVar(&commandsFlag, "c", false, "show list of commands (shorthand)")
	flag.BoolVar(&versionFlag, "version", false, "show version")
	flag.BoolVar(&versionFlag, "v", false, "show version (shorthand)")
	flag.StringVar(&outputFlag, "output", "terminal", "output format (terminal | json)")
	flag.StringVar(&outputFlag, "o", "terminal", "output format (terminal | json) (shorthand)")
}

func main() {
	flag.Parse()
	var r renderer.Renderer
	r = &renderer.TerminalRenderer{}
	ds := &datasource.BoltDataSource{}
	folder, err := osext.ExecutableFolder()
	if err != nil {
		fatalError(r, fmt.Errorf("error reading folder containing the calories binary, %v", err))
	}
	var configFile = filepath.Join(folder, configFile)
	config, err := ioutil.ReadFile(configFile)
	if err != nil {
		asciilogo()
		if os.IsNotExist(err) {
			configToSet, configErr := createConfig(config, folder, defaultDBFile)
			if configErr != nil {
				fatalError(r, configErr)
			}
			config = []byte(configToSet)
		} else {
			fatalError(r, fmt.Errorf("error reading config file at %s, %v", configFile, err))
		}
	}
	dbString := string(config)
	close, err := ds.Setup(dbString)
	if err != nil {
		fatalError(r, fmt.Errorf("could not connect to database at %s, if you want to set a new database file, please use the 'calories db' command, %v", dbString, err))
	}
	defer func() {
		closeErr := close()
		if closeErr != nil {
			fatalError(r, closeErr)
		}
	}()

	if len(flag.Args()) > 0 {
		res, err := handleSubCommand(commandFlag, commandOutputFlag, ds, r, os.Args)
		if err != nil {
			fatalError(r, err)
		}
		fmt.Fprintln(color.Output, res)
	} else {
		res, err := handleNoSubCommand(commandsFlag, outputFlag, ds, r, os.Args)
		if err != nil {
			fatalError(r, err)
		}
		fmt.Fprintln(color.Output, res)
	}
}

// handleSubCommand handles calls to subcommands
func handleSubCommand(commandFlag *flag.FlagSet, commandOutputFlag string, ds datasource.DataSource, r renderer.Renderer, args []string) (string, error) {
	err := commandFlag.Parse(args[2:])
	if err != nil {
		return "", err
	}
	command := args[1]
	if commandOutputFlag == "json" {
		r = &renderer.JSONRenderer{}
	}
	return executeCommand(ds, r, command, commandFlag.Args())
}

// handleNoSubCommand handles calls without a subcommand
func handleNoSubCommand(commandsFlag bool, outputFlag string, ds datasource.DataSource, r renderer.Renderer, args []string) (string, error) {
	if commandsFlag {
		printCommands()
		return "", nil
	}
	if versionFlag {
		return fmt.Sprintf("Current Version: %s", VERSION), nil
	}
	if outputFlag == "json" {
		r = &renderer.JSONRenderer{}
	}
	return executeCommand(ds, r, "", args)
}

// createConfig asks the user which database file to use and writes the answer into
// the .caloriesconf configuration file
func createConfig(config []byte, folder, defaultDBFile string) (string, error) {
	defaultConfig := filepath.Join(folder, defaultDBFile)
	homedir, err := homedir.Dir()
	if err != nil {
		fmt.Println("could not find home directory, using binary directory instead")
	}
	if homedir != "" {
		defaultConfig = filepath.Join(homedir, defaultDBFile)
	}
	configToSet := defaultConfig
	prompt := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("Setup database file (default: %s): ", defaultConfig)
		response, promptErr := prompt.ReadString('\n')
		if promptErr != nil {
			return "", promptErr
		}
		response = strings.TrimSpace(response)
		if response != "" {
			configToSet = response
		}
		break
	}
	configFilePath := filepath.Join(folder, configFile)
	err = ioutil.WriteFile(configFilePath, []byte(configToSet), os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("error writing config file at %s, %v", configFile, err)
	}
	fmt.Printf("config file written successfully, using database: %s\n", configToSet)
	return configToSet, nil
}

// fatalError prints the given error using the provided renderer and exits the program
func fatalError(r renderer.Renderer, fatalError error) {
	res, err := r.Error(fatalError)
	if err != nil {
		fmt.Printf("Could not print error, %v\n", err)
		os.Exit(1)
	}
	fmt.Println(res)
	os.Exit(1)
}

// executeCommand parses the subcommands and executes the associated command
// If there is no subcommand, it executes the default command
// which shows the current day/week/month
func executeCommand(ds datasource.DataSource, r renderer.Renderer, cmd string, args []string) (string, error) {
	switch cmd {
	case "weight":
		var weight string
		if len(args) > 0 {
			weight = args[0]
		}
		return checkConfig(ds, &command.WeightCommand{
			DataSource: ds,
			Renderer:   r,
			Weight:     weight,
			Mode:       len(args),
		})
	case "config":
		configCmd := command.ConfigCommand{
			DataSource: ds,
			Renderer:   r,
			Weight:     weightFlag,
			Height:     heightFlag,
			Activity:   activityFlag,
			Birthday:   birthDayFlag,
			Gender:     genderFlag,
			UnitSystem: unitFlag,
			YesMode:    yesFlag,
			Mode:       commandFlag.NFlag(),
		}
		return configCmd.Execute()
	case "add":
		var food string
		var calories string
		if len(args) >= 2 {
			calories = args[0]
			food = args[1]
		}

		return checkConfig(ds, &command.AddEntryCommand{
			DataSource: ds,
			Renderer:   r,
			Date:       dateFlag,
			Food:       food,
			Calories:   calories,
			Mode:       len(args),
		})
	case "clear":
		return checkConfig(ds, &command.ClearEntriesCommand{
			DataSource: ds,
			Renderer:   r,
			Date:       dateFlag,
			Position:   positionFlag,
			YesMode:    yesFlag,
		})
	case "export":
		return checkConfig(ds, &command.ExportCommand{
			DataSource: ds,
		})
	case "import":
		return checkConfig(ds, &command.ImportCommand{
			DataSource: ds,
			Renderer:   r,
			File:       fileFlag,
		})
	default:
		return checkConfig(ds, &command.DayCommand{
			DataSource:  ds,
			Renderer:    r,
			Week:        weekFlag,
			Month:       monthFlag,
			History:     histFlag,
			DefaultDate: defaultDateFlag,
		})
	}
}

// checkConfig checks if a config has been set and returns an error if not, explaining
// to the user how to set the config
func checkConfig(ds datasource.DataSource, cmd command.Command) (string, error) {
	_, err := ds.FetchConfig()
	if err != nil {
		return "", fmt.Errorf("no config has been set, please use: 'calories config --weight=0.0 --height=0.0 --activity=0.0 --birthday=01.01.1970 --gender=female --unit=metric' to set it")
	}
	return cmd.Execute()
}

// printCommands shows a list of available commands
func printCommands() {
	fmt.Println("You can show the HELP for each command by using")
	fmt.Println("\tCOMMAND --help")
	fmt.Println("")
	fmt.Println("You can switch the output format of each command by using")
	fmt.Println("\tCOMMAND --o=[string[terminal|json] OUTPUTFORMAT]")
	fmt.Println("")
	fmt.Println("List of Commands:")
	fmt.Println("")
	fmt.Println("- config")
	fmt.Println("\tDisplays your current configuration")
	fmt.Println("")
	fmt.Println("- config --w=[float WEIGHT] --h=[float HEIGHT] --a=[float ACTIVITY] --b=[date[dd.mm.yyyy] BIRTHDAY], --g=[string[male|female] GENDER] --u=[string[metric|imperial] UNITSYSTEM")
	fmt.Println("\tOverrides the configuration with the given values, asks for confirmation")
	fmt.Println("")
	fmt.Println("- weight")
	fmt.Println("\tDisplays your weight timeline")
	fmt.Println("")
	fmt.Println("- weight [float WEIGHT]")
	fmt.Println("\tAdds the given weight to your weight timeline with date = today")
	fmt.Println("")
	fmt.Println("- add [int CALORIES] [string FOOD]")
	fmt.Println("\tAdds an entry with the given calories and food for today")
	fmt.Println("")
	fmt.Println("- add --date=[date[dd.mm.yyyy] DATE] [int CALORIES] [string FOOD]")
	fmt.Println("\tAdds an entry with the given calories and food for the given date")
	fmt.Println("")
	fmt.Println("- clear")
	fmt.Println("\tClears the entries for the current day, asks for confirmation")
	fmt.Println("")
	fmt.Println("- clear --date=[date[dd.mm.yyyy] DATE]")
	fmt.Println("\tClears the entries for the given day, asks for confirmation")
	fmt.Println("")
	fmt.Println("- clear --position=[int POSITION]")
	fmt.Println("\tClears the entry at the given position (1-n) for the given day, asks for confirmation")
	fmt.Println("")
	fmt.Println("- export > backup.json")
	fmt.Println("\tExports the database to stdout")
	fmt.Println("")
	fmt.Println("- import --f=[string FILENAME]")
	fmt.Println("\tImports the given file to the database, overwriting all data")
}

// asciilogo prints the logo in ascii
func asciilogo() {
	fmt.Println(`
	Welcome to Calories!
	
                                                              './y/.                                
                               '........'..                  -oss+sh/'                              
                          ../+soyyssssysss+ss++o/:/-...    ':sssdmsyo.                              
                     '.-+sohyhhdddddhhdhhddddhhhhhhdhyss/--:yymhhhyy-'                              
        '.-//:-'  '.+oyhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhdddhdhhsdNNsoo-.'                              
      ./ooyyyyy+/-oyhhddhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhdhddshNNyyoohyoo/.'                          
    '/osyhhhhhhyoyyhhhyhdmmNNNNmdsyhhhhhhhhhhhhyyhhddddhhhdymNh+shhhhhysso-                         
    +yshhhsooshhhhhhyhmMMMMMMMMMMNyshhhhhhhhhhhmNMMMMMMMNdyymmNhoyssoohhyy..                        
    ohhhhy+//+ydhdydNNNNNMMMMMMMMMMmyhhhhhhydMMMMMMMMMMMMMmysmNmooo+//shhy+-                        
    ohyhhh+/ohhhhhyMy++++hMMMMMMMMMMdyhhhhydMdo+oyNMMMMMMMMNsshmmyhs+ohhhs+:                        
    -+sshhhshhhhhhmM/oo///MMMMMMMMMMdyhhhhyhN+s+//yMMMMMMMMMshydNmyhyyhhhs+.                        
     ./soyyyhhhhhhhMdso++yMMMMMMMMMNyhhhhhhsNyo++sNMMMMMMMMhsdhdmNmyhhyyoy:'                        
       .ohhhhhhhhhhhNNNNNMMMMMMMMMmshhhhhhhhsNmdNMMMMMMMMNdyhhhhhNNdyhsoy/.                         
       .+ddhhhhhhhhhhhmNNNMMMMNmhhyhhyo+oosyhymNNNNNNNNmdhhhhhhhhhmmmhhso.                          
       /shhhhhhhhhhhhhhhhhhddhhyhhhhy+//////hhhhyyhhhhhhhhhhhhhhhdhdNdhys/'                         
      ':ddhhhhhhhhhhhhhhhhddddhhhhhhdhysssydmmdhhhhhhhhhhhhhhhhhhhhdhdNdyy:                         
      -oddhhhhhhhhhhhhhhhhhhhhhhdhdNNNNNNNNNNNNmyhhhhhhhhhhhhhhhhhhhhhhhys:                         
      /shdhhhhhhhhhhhhhhhhhhhhhhhhhdmddmmmmmmmmdyhhhhhhhhhhhhhhhhhhhhhhhyo.                         
      :yhdhhhhhhhhhhhhhhhhhhhhhhhhhhhsNMdNMMyhhhhhhhhhhhhhhhhhhhhhhhhhhhdo-                         
      :sddhhhhhhhhhhhhhhhhhhhhhhhhhhhyMMhNMMyhhhhhhhhhhhhhhhhhhhhhhhhhhhh+-                         
      ./ddhhhhhhhhhhhhhhhhhhhhhhhhhhhhddydNdyhhhhhhhhhhhhhhhhhhhhhhhhhhhy+.                         
       :smhhhhhhhhhhhhhhhhhhhhhhhhhhhhhddhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhho-  .:----.'               
       '+mddhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhs/'.::::::::'              
       ./myhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhyo/:::::::::::-''           
       .+mNNNmmdmddhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh/::::::::::::/:/.          
        /yNmNNNNNNmmmddhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh/:::::::::::::::.          
 ..:://+:/dymdmNNNNNmdmdddhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhho+++syyyoyy/:+:-.          
.+dNNNmmdhhddhhhhNNmdNNNNNmddhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhs/-+sNmmmy/-:.'            
-omNNNNNNNmdhhhhddddNNdNNNNNmmmmddhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhymdmdmmmy:                 
 -/osy/+oyddhhhhhhhddhmNNdNNNNNNNNmddhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhdNNmmmmNs                  
       ''/ddhhhhhhhhhhhddmNNdNNNNNNNmdddhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhddy+yommmo                  
        -sddhhhhhhhhhhhhhhdddmmdNNNdmNNNmmddhhhhhhhhhhhhhhhhhhhhhhhhhhhhhd+.':yys-                  
        +dddhhhhhhhhhhhhhhhhhddddmdmNmmNNNNNmmmdhhhhhhhhhhhhhhhhhhhhhhhhhd/:  '-'                   
       -/ddhhhhhhhhhhhhhhhhhhhhhhhhdmdNNNmmNNNNNNmdhhdhhhhhhhhhhhhhhhhhhhho-                        
       :sdhhhhhhhhhhhhhhhhhhhhhhhhhhhdddddmNNdNNNNNNNmNhdddddddhhhhhhhhhhds/'                       
       +sdhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhdddhmmmdNNNNmdNNNNNmmmdddhhhhhhdy:-                       
       :hdhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhddddmmNdmddmdNNmNNNmmdhdhhhh/-                       
       :mdhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhdddhhhddhdmmmNNdmNmddhh:                       
       :ddhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhddhhdmdmNmNNNdho:.                     
       :mhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhdhdmdmmmdmmds:'                   
       :ddhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhdhhymymdy-                  
       /sdhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhdsssdmNms/'                
       :odhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh/yhmmNmho-               
        +ddhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhys:-shmmNmd/'              
        :odhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhy+- :oymNNmo:              
         /yhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh/'  /yhmmhh/              
         '+hddhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhho-   ':yhmmdo++:.          
          '/hmdhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhdyh/     :ohmNddNmhs:.        
           ':hmdhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhdyy/'     -/hdNmhyyhhh+'       
             :shddhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhyy+'     .+yysmNy++yyh+-'-:'   
             '/+yyhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhyo'    .+hhdyhNhyhyhs+/osso:. 
            /shhdmddyhhdhdhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhdhhdmNNNdy+.   :shmmNyydddho+ydddsy+/ 
          -+hmNNNNNmyoshhhhhdddhhhhhhhhhhhhhhhhhhhhhhhhhhdddhhsshmNNmhmmo:   :/hmmmymssoyhddddyho:  
         -ohmmNNNNdh/.'.::+++hhhhhhddhhhddhhhhhhhhhhhdddddhs+-./ohmNNdyy:.    /ydmNmyyyddNhdhyo/.   
         ./dddNNNds:'        ../:/+sohhyhhyhhhhhdddhhhyoo/.'    '//ssys+'     ':ydhmdmdmhdhs/-'     
          :ohhmhyo-                 '....//o+yyo:-::-.-'          '''..        '/syysysso/..        
           '.:/..                                                                '.-..-.            
                                                                                                    
                                                                                                    

	`)
}
