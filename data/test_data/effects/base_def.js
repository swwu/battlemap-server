/*
 * Implement the base character rules of the Pathfinder system -
 *
 * stats, stat bonuses, etc
 *
 * In general things that "define" variables should be priority 0, and things
 * that in
 */

var ability_scores = ["str","dex","con","int","wis","cha"];

var bonus_types = ["alchemical", "armor", "armor_enhancement", "circumstance",
    "competence", "deflection", "dodge", "enhancement", "inherent", "insight",
    "luck", "morale", "natural_armor", "natural_armor_enhancement", "profane",
    "racial", "resistance", "sacred", "shield", "shield_enhancement", "size",
    "trait", "untyped"];
var touch_ac_exclude = ["armor", "armor_enhancement", "natural_armor",
    "natural_armor_enhancement", "shield", "shield_enhancement"];
var flatfooted_ac_exclude = ["dodge"];

var str_skills = ["skill_climb", "skill_swim"];
var dex_skills = ["skill_acrobatics", "skill_disabledevice",
    "skill_escapeartist", "skill_fly", "skill_ride", "skill_sleightofhand",
    "skill_stealth"];
var con_skills = [];
var int_skills = ["skill_appraise", "skill_craft", "skill_knowledge",
    "skill_linguistics", "skill_spellcraft"];
var wis_skills = ["skill_heal", "skill_perception", "skill_profession",
    "skill_sensemotive", "skill_survival"];
var cha_skills = ["skill_bluff", "skill_diplomacy", "skill_disguise",
    "skill_handleanimal", "skill_intimidate", "skill_perform",
    "skill_usemagicdevice"];

// returns accumulators for all the bonuses 
var generateBonusNames = function(stat_name, exclude) {
  return bonus_types.map(function(bonus_type) {
    return stat_name + "_" + bonus_type + "_bonus";
  }).filter(function(bonus_type) {
    return !(exclude && exclude.indexOf(bonus_type) > -1);
  });
};

var generateBonusAccums = function(entity, stat_name) {
  generateBonusNames(stat_name).forEach(function(stat_name,index) {
    // remember that some bonuses stack, others don't
    entity.vars.newAccum({
      id: stat_name,
      op: ["circumstance","dodge","racial","untyped"].indexOf(bonus_types[index]) > -1 ? '+' : 'max',
      init: 0
    });
  });
};

var sumBonuses = function(deps, stat_name, exclude) {
  return generateBonusNames(stat_name,exclude).reduce(function(accum, stat_name) {
    return accum + deps[stat_name];
  },0);
}


define.effect({
  id: "baseStats",
  displayName: "Base Stats",
  displayType: "base",

  onEffect: function(entity) {
    entity.vars.new({
      id:"base_stats_proxy",
      modifies: ability_scores.map(function(score){return score+"_base";}),
      onEval: function(deps, mods) {
        for (var i=0; i<ability_scores.length; i++) {
          mods[ability_scores[i]+"_base"](14);
        }
      }
    })
  }
});

define.effect({
  id: "baseEntityRules",
  displayName: "Base Rule Calculations",
  displayType: "base",

  onEffect: function(entity) {
    // ability scores
    for (var i=0; i<ability_scores.length; i++) {
      // wrap the scope in a function so that callbacks reference the correct
      // value of i
      (function (score) {
        entity.vars.newAccum({
          id: score+"_base",
          op: "+",
          init: 0
        })
        generateBonusAccums(entity,score);
        entity.vars.new({
          id: score,
          depends: [score+"_base"].concat(generateBonusNames(score)),
          onEval: function(deps) {
            return deps[score+"_base"] + sumBonuses(deps, score);
          }
        });
        entity.vars.new({
          id: score+"_mod",
          depends: [score],
          onEval: function(deps) {
            return Math.floor((deps[score]-10)/2);
          }
        });

      })(ability_scores[i]);
    };

    // hp
    entity.vars.newAccum({
      id: "hp",
      op: "+",
      init: 0
    });

    // calculate BAB/melee AB/ranged AB/CMB/CMD
    entity.vars.newAccum({
      id: "bab",
      op: "+",
      init: 0
    });
    generateBonusAccums(entity,"melee_ab");
    entity.vars.newAccum({
      id: "melee_ab_stat_bonus",
      op: "max",
      init: 0
    })
    entity.vars.new({
      id: "melee_ab",
      depends: ["bab", "melee_ab_stat_bonus"].concat(generateBonusNames("melee_ab")),
      onEval: function(deps) {
        return deps.bab + deps.melee_ab_stat_bonus + sumBonuses(deps, "melee_ab");
      }
    });
    generateBonusAccums(entity,"ranged_ab");
    entity.vars.newAccum({
      id: "ranged_ab_stat_bonus",
      op: "max",
      init: 0
    })
    entity.vars.new({
      id: "ranged_ab",
      depends: ["bab", "ranged_ab_stat_bonus"].concat(generateBonusNames("ranged_ab")),
      onEval: function(deps) {
        return deps.bab + deps.ranged_ab_stat_bonus + sumBonuses(deps, "ranged_ab");
      }
    });
    entity.vars.new({
      id: "ab_proxy",
      depends: ["str_mod", "dex_mod"],
      modifies: ["melee_ab_stat_bonus", "ranged_ab_stat_bonus"],
      onEval: function(deps, mods) {
        mods.melee_ab_stat_bonus(deps.str_mod);
        mods.ranged_ab_stat_bonus(deps.dex_mod);
      }
    })

    // ac
    entity.vars.newAccum({
      id: "ac_base",
      op: "+",
      init: 10
    });
    entity.vars.newAccum({
      id: "ac_stat_bonus",
      op: "max",
      init: 0
    })
    generateBonusAccums(entity,"ac");
    entity.vars.new({
      id: "ac",
      depends: ["ac_base","ac_stat_bonus"].concat(generateBonusNames("ac")),
      onEval: function(deps) {
        return deps.ac_base + deps.ac_stat_bonus + sumBonuses(deps, "ac");
      }
    })
    entity.vars.new({
      id: "ac_touch",
      depends: ["ac_base","ac_stat_bonus"].concat(generateBonusNames("ac",touch_ac_exclude)),
      onEval: function(deps) {
        return deps.ac_base + deps.ac_stat_bonus + sumBonuses(deps, "ac", touch_ac_exclude);
      }
    })
    entity.vars.new({
      id: "ac_flatfooted",
      depends: ["ac_base"].concat(generateBonusNames("ac",flatfooted_ac_exclude)),
      onEval: function(deps) {
        return deps.ac_base + sumBonuses(deps, "ac", flatfooted_ac_exclude);
      }
    })
    entity.vars.new({
      id: "ac_proxy",
      depends: ["dex_mod"],
      modifies: ["ac_stat_bonus"],
      onEval: function(deps, mods) {
        mods.ac_stat_bonus(deps.dex_mod);
      }
    })





    // saves
    entity.vars.newAccum({
      id: "will_save_base",
      op: "+",
      init: 0
    });
    generateBonusAccums(entity,"will_save");
    entity.vars.new({
      id: "will_save",
      depends: ["will_save_base","wis_mod"].concat(generateBonusNames("will_save")),
      onEval: function(deps) {
        return deps.will_save_base + deps.wis_mod + sumBonuses(deps, "will_save");
      }
    })
    entity.vars.newAccum({
      id: "fort_save_base",
      op: "+",
      init: 0
    });
    generateBonusAccums(entity,"fort_save");
    entity.vars.new({
      id: "fort_save",
      depends: ["fort_save_base","con_mod"].concat(generateBonusNames("fort_save")),
      onEval: function(deps) {
        return deps.fort_save_base + deps.con_mod + sumBonuses(deps, "fort_save");
      }
    })
    entity.vars.newAccum({
      id: "ref_save_base",
      op: "+",
      init: 0
    });
    generateBonusAccums(entity,"ref_save");
    entity.vars.new({
      id: "ref_save",
      depends: ["ref_save_base","dex_mod"].concat(generateBonusNames("ref_save")),
      onEval: function(deps) {
        return deps.ref_save_base + deps.dex_mod + sumBonuses(deps, "ref_save");
      }
    })

    // testing stuff
    entity.vars.new({
      id: "test_proxy_1",
      modifies: ["fighter_lvl","will_save_insight_bonus","will_save_untyped_bonus"],
      onEval: function(deps,mods) {
        mods.will_save_insight_bonus(1);
        mods.will_save_insight_bonus(2);
        mods.will_save_untyped_bonus(1);
        mods.will_save_untyped_bonus(2);
        mods.fighter_lvl(10);
      }
    })
  }
})
