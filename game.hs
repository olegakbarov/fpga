import Control.Applicative
import System.Random

data Suit = Hearts
          | Spades
          | Clubs
          | Diamonds
          deriving (Eq, Enum, Ord, Read, Show)

data Rank = Ace
          | Two
          | Three
          | Four
          | Five
          | Six
          | Seven
          | Eight
          | Nine
          | Ten
          | Jack
          | Queen
          | King
          deriving (Eq, Ord, Bounded, Enum, Show, Read)

data Card = Card Rank Suit deriving (Eq, Show, Read)

createDeck f [] = [Card r s | r <- [Ace .. King], s <- [Hearts .. Diamonds]]
-- data Deck = createDeck

data Hand = Hand [Card]


