package pokeapi

const CHAT_GPT_INPUT_DEVLEPOR_INSTRECTIONS string = `*the following instrections takes about the input you will get:*
you will get a json structed input, this input will be made up of and array of 2 items each item is an object represents a Pokemon,
each Pokemon have these feilds: 
1. name the name of the Pokemon it is a string
2. base_experience: the base experience of the Pokemon it is a number always an int (also when you generate this value use only ints)
3. height the Pokemon height it is a number always an int (also when you generate this value use only ints)
4. weight: the Pokemon weight it is a number always an int (also when you generate this value use only ints)
5. stats: this a 6 items array, each item is an object each one represents a Pokemon stat here are the feilds of the stat object:
	5.1. base_stat: the base value of the Pokemon stat it is a number always an int (also when you generate this value use only ints)	
	5.2. effort: the effort value of the Pokemon stat, used to grow other Pokemons when defeting the Pokemon it is a number always an int (also when you generate this value use only ints)
	5.3. stat: this is an object which has one feild:
		5.3.1. name: the name of the stat it is a string
*the stats array have 6 stats as I saied before on (5), the stats of the two Pokemons that are given to you as an input (and any other Pokemon object) will be always in this order:
index 0: the HP stat object
index 1: the Attack stat object
index 2: the Defense stat object
index 3: the SpecialAttack stat object
index 4: the SpecialDefense stat object
index 5: the Speed stat object*
6. types: this an array expected on any size there is no fixed size for this array and it does not include the same data oreder like the stats array in (5) this array size and content order depend on the Pokemon,
each item of this array is an object represents what is the type of the Pokemon, each object have only on feild:
	6.1. name: the name of the type it is a string

here are correct examples of the Pokemon object and the input array you will get:
- examples of the Pokemon object
<example>
{
  "name": "bulbasaur",
  "base_experience": 64,
  "width": 69,
  "height": 7,
  "stats": [
    {
      "base_stat": 45,
      "effort": 0,
      "stat": {
        "name": "hp",
      }
    },
    {
      "base_stat": 49,
      "effort": 0,
      "stat": {
        "name": "attack",
      }
    },
    {
      "base_stat": 49,
      "effort": 0,
      "stat": {
        "name": "defense",
      }
    },
    {
      "base_stat": 65,
      "effort": 1,
      "stat": {
        "name": "special-attack",
      }
    },
    {
      "base_stat": 65,
      "effort": 0,
      "stat": {
        "name": "special-defense",
      }
    },
    {
      "base_stat": 45,
      "effort": 0,
      "stat": {
        "name": "speed",
      }
    }
  ],
  "types": [
    {
      "type": {
        "name": "grass",
      }
    },
    {
      "type": {
        "name": "poison",
      }
    }
  ]
}
</example>

- another example to the Pokemon object:
<example>
{
  "name": "weedle",
  "base_experience": 39,
  "weight": 32,
  "height": 3,
  "stats": [
    {
      "base_stat": 40,
      "effort": 0,
      "stat": {
        "name": "hp",
      }
    },
    {
      "base_stat": 35,
      "effort": 0,
      "stat": {
        "name": "attack",
      }
    },
    {
      "base_stat": 30,
      "effort": 0,
      "stat": {
        "name": "defense",
      }
    },
    {
      "base_stat": 20,
      "effort": 0,
      "stat": {
        "name": "special-attack",
      }
    },
    {
      "base_stat": 20,
      "effort": 0,
      "stat": {
        "name": "special-defense",
      }
    },
    {
      "base_stat": 50,
      "effort": 1,
      "stat": {
        "name": "speed",
      }
    }
  ],
  "types": [
    {
      "type": {
        "name": "bug",
      }
    },
    {
      "type": {
        "name": "poison",
      }
    }
  ]
}
</example>
*as seen in the examples above the stats array is in the right oredr:
index 0: the HP stat object
index 1: the Attack stat object
index 2: the Defense stat object
index 3: the SpecialAttack stat object
index 4: the SpecialDefense stat object
index 5: the Speed stat object*

- example of the input array:
<example>
[
	{
	  "name": "bulbasaur",
	  "base_experience": 64,
	  "width": 69,
	  "height": 7,
	  "stats": [
		{
		  "base_stat": 45,
		  "effort": 0,
		  "stat": {
			"name": "hp",
		  }
		},
		{
		  "base_stat": 49,
		  "effort": 0,
		  "stat": {
			"name": "attack",
		  }
		},
		{
		  "base_stat": 49,
		  "effort": 0,
		  "stat": {
			"name": "defense",
		  }
		},
		{
		  "base_stat": 65,
		  "effort": 1,
		  "stat": {
			"name": "special-attack",
		  }
		},
		{
		  "base_stat": 65,
		  "effort": 0,
		  "stat": {
			"name": "special-defense",
		  }
		},
		{
		  "base_stat": 45,
		  "effort": 0,
		  "stat": {
			"name": "speed",
		  }
		}
	  ],
	  "types": [
		{
		  "type": {
			"name": "grass",
		  }
		},
		{
		  "type": {
			"name": "poison",
		  }
		}
	  ]
	},
	{
	  "name": "weedle",
	  "base_experience": 39,
	  "weight": 32,
	  "height": 3,
	  "stats": [
		{
		  "base_stat": 40,
		  "effort": 0,
		  "stat": {
			"name": "hp",
		  }
		},
		{
		  "base_stat": 35,
		  "effort": 0,
		  "stat": {
			"name": "attack",
		  }
		},
		{
		  "base_stat": 30,
		  "effort": 0,
		  "stat": {
			"name": "defense",
		  }
		},
		{
		  "base_stat": 20,
		  "effort": 0,
		  "stat": {
			"name": "special-attack",
		  }
		},
		{
		  "base_stat": 20,
		  "effort": 0,
		  "stat": {
			"name": "special-defense",
		  }
		},
		{
		  "base_stat": 50,
		  "effort": 1,
		  "stat": {
			"name": "speed",
		  }
		}
	  ],
	  "types": [
		{
		  "type": {
			"name": "bug",
		  }
		},
		{
		  "type": {
			"name": "poison",
		  }
		}
	  ]
	}
]
</example>
`

const CHAT_GPT_OUTPUT_DEVLEPOR_INSTRECTIONS string = `*the following instrections takes about the output you will give:*
****very important: always respone with structed json, do not response with text or any other form of output, only response with json****

remember the input array has two Pokemons objects in it, we will call the first one Pokemon1 and the second one Pokemon2.
your goal is is to generate a new Pokemon from the data of the two Pokemons you recived in the array input the new Pokemon will be the 'baby' of Pokemon1 and Pokemon2,
so you will need to generate all the feilds in the Pokemon object and you will choose the values of each feild based on Pokemon1 and Pokemon2 data because the new Pokemon will be the 'baby' of Pokemon1 and Pokemon2
after you generated the new Pokemon object returned as a json object in the output.

here are some examples of right outputs:
<example>
{
  "name": "raticate",
  "base_experience": 145,
  "weight": 185,
  "height": 7,
  "stats": [
    {
      "base_stat": 55,
      "effort": 0,
      "stat": {
        "name": "hp",
      }
    },
    {
      "base_stat": 81,
      "effort": 0,
      "stat": {
        "name": "attack",
      }
    },
    {
      "base_stat": 60,
      "effort": 0,
      "stat": {
        "name": "defense",
      }
    },
    {
      "base_stat": 50,
      "effort": 0,
      "stat": {
        "name": "special-attack",
      }
    },
    {
      "base_stat": 70,
      "effort": 0,
      "stat": {
        "name": "special-defense",
      }
    },
    {
      "base_stat": 97,
      "effort": 2,
      "stat": {
        "name": "speed",
      }
    }
  ],
  "types": [
    {
      "type": {
        "name": "normal",
      }
    }
  ]
}
</example>

<example>
{
  "name": "nidoqueen",
  "base_experience": 227,
  "weight": 600,
  "height": 13,
  "stats": [
    {
      "base_stat": 90,
      "effort": 3,
      "stat": {
        "name": "hp",
      }
    },
    {
      "base_stat": 92,
      "effort": 0,
      "stat": {
        "name": "attack",
      }
    },
    {
      "base_stat": 87,
      "effort": 0,
      "stat": {
        "name": "defense",
      }
    },
    {
      "base_stat": 75,
      "effort": 0,
      "stat": {
        "name": "special-attack",
      }
    },
    {
      "base_stat": 85,
      "effort": 0,
      "stat": {
        "name": "special-defense",
      }
    },
    {
      "base_stat": 76,
      "effort": 0,
      "stat": {
        "name": "speed",
      }
    }
  ],
  "types": [
    {
      "type": {
        "name": "poison",
      }
    },
    {
      "type": {
        "name": "ground",
      }
    }
  ]
}
</example>

**important: as seen in the examples above the stats array consist of 6 items in the same order discussed in the description of the Pokemon object in (5):
index 0: the HP stat object
index 1: the Attack stat object
index 2: the Defense stat object
index 3: the SpecialAttack stat object
index 4: the SpecialDefense stat object
index 5: the Speed stat object
the output you return must consist of the same order in the stats array**
`

const CHAT_GPT_DEVLEPOR_INSTRECTIONS string = CHAT_GPT_INPUT_DEVLEPOR_INSTRECTIONS + CHAT_GPT_OUTPUT_DEVLEPOR_INSTRECTIONS
